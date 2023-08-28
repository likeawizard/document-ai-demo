package dispatch

import (
	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/expensebot"
)

type EventMsg struct {
	record database.Record
	msg    string
}

type EventChan struct {
	c chan EventMsg
}

type Dispatch struct {
	processor expensebot.DocumentProcessor
	pipe      map[uuid.UUID]Pipeline
}

type Pipeline struct {
	r    database.Record
	done bool
}

func (d *Dispatch) New(record database.Record) {
	// d.pipe = append(d.pipe, Pipeline{r: record, done: false})
}

func (ec EventChan) NewRecord(record database.Record) {
	ec.c <- EventMsg{record: record, msg: "new"}
}

func (ec EventChan) ProcessStarted(record database.Record) {
	ec.c <- EventMsg{record: record, msg: "processed"}
}

func (ec EventChan) ProcesseDone(record database.Record) {
	ec.c <- EventMsg{record: record, msg: "new"}
}
