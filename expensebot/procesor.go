package expensebot

import (
	"errors"
	"fmt"
	"log"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
)

var Processor DocumentProcessor

type DocumentProcessor interface {
	Process(record database.Record) error
	Schema() string
}

func NewDocumentProcessor(cfg config.ProcessorCfg) (DocumentProcessor, error) {
	if cfg == nil {
		return nil, errors.New("no config provided to document processor")
	}
	switch v := cfg.(type) {
	case *config.DocuIntelCfg:
		return NewDocuIntel(*v), nil
	case *config.DocumentAICfg:
		return NewGoogleDocumentAI(*v), nil
	default:
		return nil, fmt.Errorf("unsupported processor driver: %s", v.Driver())
	}
}

// TODO you do not belong here
func updateWithStatus(record database.Record, status database.Status) {
	if config.App.Debug {
		log.Printf("record status updated for uuid: %s from: %s to: %s\n", record.Id, record.Status, status)
	}
	record.Status = status
	err := database.Instance.Update(record)
	if err != nil {
		log.Printf("failed updating record: %v\n", err)
	}
}
