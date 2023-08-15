package expensebot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/store"
	"google.golang.org/api/option"
)

var Processor DocumentProcessor

type DocumentProcessor interface {
	Process(record database.Record) error
}

type GoogleDocumentAI struct {
	credsFile string
	endpoint  string
	name      string
}

func NewDocumentProcessor(cfg config.DocumentAICfg) DocumentProcessor {
	return NewGoogleDocumentAI(cfg)
}

func NewGoogleDocumentAI(cfg config.DocumentAICfg) *GoogleDocumentAI {
	return &GoogleDocumentAI{
		credsFile: cfg.CredsFile,
		endpoint:  fmt.Sprintf("%s-documentai.googleapis.com:443", cfg.Location),
		name:      fmt.Sprintf("projects/%s/locations/%s/processors/%s", cfg.ProjectId, cfg.Location, cfg.ProcessorId),
	}
}

func (docAI *GoogleDocumentAI) Process(record database.Record) error {
	ctx := context.Background()
	client, err := docAI.newDocumentProcessorClient(ctx)
	if err != nil {
		return err
	}
	req, err := docAI.newDocumentProcessorRequest(ctx, record)
	if err != nil {
		return err
	}

	go func(record database.Record) {
		resp, err := client.ProcessDocument(ctx, req)
		if err != nil {
			log.Print("failed GoogleDocumentAI ProcessDocument call:", err)
			updateWithStatus(record, database.S_FAILED)
			return
		}
		doc := resp.GetDocument()
		json, err := json.Marshal(doc)
		if err != nil {
			log.Print("failed GoogleDocumentAI doc to json marshal:", err)
			updateWithStatus(record, database.S_FAILED)
			return
		}
		jsonPath := fmt.Sprintf("%s.json", record.Id)
		err = store.File.Store(jsonPath, bytes.NewReader(json))
		if err != nil {
			log.Print("failed GoogleDocumentAI file storage:", err)
			updateWithStatus(record, database.S_FAILED)
			return
		}
		record.JSON = jsonPath
		updateWithStatus(record, database.S_READY)
	}(record)
	return nil
}

func updateWithStatus(record database.Record, status database.Status) {
	if config.App.Debug {
		log.Printf("record status updated for uuid: %s from: %s to: %s\n", record.Id, record.Status, status)
	}
	record.Status = status
	database.Instance.Update(record)
}

func (docAI *GoogleDocumentAI) newDocumentProcessorClient(ctx context.Context) (*documentai.DocumentProcessorClient, error) {
	auth := option.WithCredentialsFile(docAI.credsFile)
	endpointOpt := option.WithEndpoint(docAI.endpoint)
	return documentai.NewDocumentProcessorClient(ctx, endpointOpt, auth)
}

func (docAI *GoogleDocumentAI) newDocumentProcessorRequest(ctx context.Context, record database.Record) (*documentaipb.ProcessRequest, error) {
	f, err := store.File.Get(record.Path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	req := documentaipb.ProcessRequest{
		SkipHumanReview: false,
		Name:            docAI.name,
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  data,
				MimeType: record.MimeType,
			},
		},
	}
	return &req, nil

}
