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
	Receipt database.Receipt
	Msg     string
	Data    map[string]string
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
		log.Printf("New event for %s with msg %s data : '%+v'", event.Receipt.Id, event.Msg, event.Data)
		switch event.Msg {
		case msgNew:
			go pe.DispatchProcess(event.Receipt)
		case msgProcessed:
			go pe.DispatchDataTransform(event.Receipt, event.Data["schema"])
		case msgTransformed:
			go pe.DispatchPostProcess(event.Receipt)
		case msgDone:
			go pe.DispatchDone(event.Receipt)
		case msgFailed:
			go pe.DispatchFailed(event.Receipt, errors.New(event.Data["err"]))
		default:
			err := fmt.Errorf("unknown event message: '%s'", event.Msg)
			go pe.DispatchFailed(event.Receipt, err)
		}
	}
}

func (pe *ExpenseEngine) DispatchProcess(receipt database.Receipt) {
	err := pe.processService.Process(receipt)
	if err != nil {
		pe.eventChan.MsgFailed(receipt, err)
		return
	}
	receipt.Status = msgProcessed
	pe.Db.Update(receipt)
	pe.eventChan.MsgProcessed(receipt, pe.processService.Processor.Schema())
}

func (pe *ExpenseEngine) DispatchDataTransform(receipt database.Receipt, schema string) {
	err := pe.transformService.Transform(receipt, schema)
	if err != nil {
		pe.eventChan.MsgFailed(receipt, err)
		return
	}
	receipt.Status = msgTransformed
	pe.Db.Update(receipt)
	pe.eventChan.MsgTransformed(receipt)
}

func (pe *ExpenseEngine) DispatchPostProcess(receipt database.Receipt) {
	file, err := pe.postProcessService.FileStore.Get(receipt.GetJsonPath())
	if err != nil {
		pe.eventChan.MsgFailed(receipt, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		pe.eventChan.MsgFailed(receipt, err)
		return
	}

	var exp transform.Expense
	err = json.Unmarshal(data, &exp)
	if err != nil {
		pe.eventChan.MsgFailed(receipt, err)
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
		pe.eventChan.MsgFailed(receipt, err)
		return
	}

	err = pe.postProcessService.FileStore.Store(receipt.GetExpensePath(), bytes.NewReader(data))
	if err != nil {
		pe.eventChan.MsgFailed(receipt, err)
		return
	}

	pe.eventChan.MsgDone(receipt)
}

func (pe *ExpenseEngine) DispatchDone(receipt database.Receipt) {
	receipt.Status = msgDone
	pe.Db.Update(receipt)
}

func (pe *ExpenseEngine) DispatchFailed(receipt database.Receipt, err error) {
	log.Printf("process pipeline failed: %s", err)
	receipt.Status = msgFailed
	pe.Db.Update(receipt)
}

func (ec EventChan) MsgNew(receipt database.Receipt) {
	ec <- EventMsg{Receipt: receipt, Msg: msgNew}
}

func (ec EventChan) MsgProcessed(receipt database.Receipt, schema string) {
	ec <- EventMsg{Receipt: receipt, Msg: msgProcessed, Data: map[string]string{"schema": schema}}
}

func (ec EventChan) MsgTransformed(receipt database.Receipt) {
	ec <- EventMsg{Receipt: receipt, Msg: msgTransformed}
}

func (ec EventChan) MsgDone(receipt database.Receipt) {
	ec <- EventMsg{Receipt: receipt, Msg: msgDone}
}

func (ec EventChan) MsgFailed(receipt database.Receipt, err error) {
	ec <- EventMsg{Receipt: receipt, Msg: msgFailed, Data: map[string]string{"err": err.Error()}}
}
