package expense

import (
	"fmt"
	"log"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
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
	eventChan        EventChan
	processService   *processor.ProcessorServcie
	transformService *transform.DataTransformService
	Db               database.DB
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

	return &pe, nil
}

func (pe *ExpenseEngine) GetSendChan() EventChan {
	return pe.eventChan
}

func (pe *ExpenseEngine) Listen() {
	for event := range pe.eventChan {
		switch event.Msg {
		case msgNew:
			go pe.DispatchProcess(event.Record)
		case msgProcessed:
			go pe.DispatchDataTransform(event.Record, event.Data["schema"])
		case msgTransformed:
			// TODO initiate postProcess - convert currency and translate
		case msgDone:
			go pe.DispatchDone(event.Record)
		case msgFailed:
			go pe.DispatchFailed(event.Record)
		default:
			log.Printf("unknown event message: '%s'", event.Msg)
			go pe.DispatchFailed(event.Record)
		}
	}
}

func (pe *ExpenseEngine) DispatchProcess(record database.Record) {
	err := pe.processService.Process(record)
	if err != nil {
		pe.eventChan.MsgFailed(record)
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
		pe.eventChan.MsgFailed(record)
		return
	}
	record.Status = msgTransformed
	pe.Db.Update(record)
	pe.eventChan.MsgTransformed(record)
}

func (pe *ExpenseEngine) DispatchDone(record database.Record) {
	record.Status = msgDone
	pe.Db.Update(record)
}

func (pe *ExpenseEngine) DispatchFailed(record database.Record) {
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

func (ec EventChan) MsgFailed(record database.Record) {
	ec <- EventMsg{Record: record, Msg: msgFailed}
}
