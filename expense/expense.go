package expense

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/postprocess"
	"github.com/likeawizard/document-ai-demo/processor"
	"github.com/likeawizard/document-ai-demo/transform"
)

type EventChan chan EventMsg

type EventMsg struct {
	Record database.Record
	Msg    string
	Data   map[string]string
}

type ExpenseEngine struct {
	eventChan          EventChan
	processService     *processor.ProcessorServcie
	transformService   *transform.DataTransformService
	postProcessService *postprocess.PostProcessService
	Db                 database.DB
}

const (
	msgNew         = "new"
	msgProcessed   = "processed"
	msgTransformed = "transformed"
	msgDone        = "done"
	msgFailed      = "failed"
)

func NewExpenseEngine(cfg config.Config) (*ExpenseEngine, error) {
	pe := ExpenseEngine{
		eventChan: make(chan EventMsg),
	}
	ps, err := processor.NewProcessorService(cfg)
	if err != nil {
		return nil, err
	}
	pe.processService = ps

	dts, err := transform.NewDataTransformService(cfg)
	if err != nil {
		return nil, err
	}
	pe.transformService = dts

	db, err := database.NewDataBase(cfg.Db)
	if err != nil {
		return nil, err
	}
	pe.Db = db

	pps, err := postprocess.NewPostProcessService(cfg)
	if err != nil {
		return nil, err
	}
	pe.postProcessService = pps

	return &pe, nil
}

func (pe *ExpenseEngine) GetSendChan() EventChan {
	return pe.eventChan
}

func (pe *ExpenseEngine) Listen() {
	for event := range pe.eventChan {
		log.Printf("New event for %s with msg %s data : '%+v'", event.Record.Id, event.Msg, event.Data)
		switch event.Msg {
		case msgNew:
			go pe.DispatchProcess(event.Record)
		case msgProcessed:
			go pe.DispatchDataTransform(event.Record, event.Data["schema"])
		case msgTransformed:
			go pe.DispatchPostProcess(event.Record)
		case msgDone:
			go pe.DispatchDone(event.Record)
		case msgFailed:
			go pe.DispatchFailed(event.Record, errors.New(event.Data["err"]))
		default:
			err := fmt.Errorf("unknown event message: '%s'", event.Msg)
			go pe.DispatchFailed(event.Record, err)
		}
	}
}

func (pe *ExpenseEngine) DispatchProcess(record database.Record) {
	err := pe.processService.Process(record)
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}
	record.Status = msgProcessed
	record.JSON = fmt.Sprintf("%s.json", record.Id)
	pe.Db.Update(record)
	pe.eventChan.MsgProcessed(record, pe.processService.Processor.Schema())
}

func (pe *ExpenseEngine) DispatchDataTransform(record database.Record, schema string) {
	err := pe.transformService.Transform(record, schema)
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}
	record.Status = msgTransformed
	pe.Db.Update(record)
	pe.eventChan.MsgTransformed(record)
}

func (pe *ExpenseEngine) DispatchPostProcess(record database.Record) {
	file, err := pe.postProcessService.FileStore.Get(record.JSON)
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}

	var exp transform.Expense
	err = json.Unmarshal(data, &exp)
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}

	cpp, err := pe.postProcessService.GetCurrencyPostProcess(exp)
	if err == nil {
		err = pe.postProcessService.CurrencyService.GetConversionRate(cpp)
		if err == nil {
			cpp.Apply(&exp)
		}
	}

	tpp, err := pe.postProcessService.GetTranslationPostProcess(exp)
	if err == nil {
		err = pe.postProcessService.TranslationService.Translate(tpp)
		if err == nil {
			tpp.Apply(&exp)
		}
	}

	data, err = json.Marshal(exp)
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}

	err = pe.postProcessService.FileStore.Store(record.JSON, bytes.NewReader(data))
	if err != nil {
		pe.eventChan.MsgFailed(record, err)
		return
	}

	pe.eventChan.MsgDone(record)
}

func (pe *ExpenseEngine) DispatchDone(record database.Record) {
	record.Status = msgDone
	pe.Db.Update(record)
}

func (pe *ExpenseEngine) DispatchFailed(record database.Record, err error) {
	log.Printf("process pipeline failed: %s", err)
	record.Status = msgFailed
	pe.Db.Update(record)
}

func (ec EventChan) MsgNew(record database.Record) {
	ec <- EventMsg{Record: record, Msg: msgNew}
}

func (ec EventChan) MsgProcessed(record database.Record, schema string) {
	ec <- EventMsg{Record: record, Msg: msgProcessed, Data: map[string]string{"schema": schema}}
}

func (ec EventChan) MsgTransformed(record database.Record) {
	ec <- EventMsg{Record: record, Msg: msgTransformed}
}

func (ec EventChan) MsgDone(record database.Record) {
	ec <- EventMsg{Record: record, Msg: msgDone}
}

func (ec EventChan) MsgFailed(record database.Record, err error) {
	ec <- EventMsg{Record: record, Msg: msgFailed, Data: map[string]string{"err": err.Error()}}
}
