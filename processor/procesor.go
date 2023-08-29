package processor

import (
	"errors"
	"fmt"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/store"
)

var Processor DocumentProcessor

type ProcessorServcie struct {
	Processor DocumentProcessor
	FileStore store.FileStore
	Db        database.DB
}

type DocumentProcessor interface {
	Process(receipt database.Receipt, fs store.FileStore) error
	Schema() string
}

func NewProcessorService(cfg config.Config) (*ProcessorServcie, error) {
	ps := ProcessorServcie{}

	processor, err := NewDocumentProcessor(cfg.Processor)
	if err != nil {
		return nil, err
	}
	ps.Processor = processor

	db, err := database.NewDataBase(cfg.Db)
	if err != nil {
		return nil, err
	}
	ps.Db = db

	store, err := store.NewFileStore(cfg.Store)
	if err != nil {
		return nil, err
	}
	ps.FileStore = store

	return &ps, nil
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

func (ps *ProcessorServcie) Process(receipt database.Receipt) error {
	err := ps.Processor.Process(receipt, ps.FileStore)
	if err != nil {
		return err
	}
	return nil
}
