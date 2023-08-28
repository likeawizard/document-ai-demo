package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/store"
	"google.golang.org/api/option"
)

type GoogleDocumentAI struct {
	credsFile string
	endpoint  string
	name      string
}

func NewGoogleDocumentAI(cfg config.DocumentAICfg) *GoogleDocumentAI {
	return &GoogleDocumentAI{
		credsFile: cfg.CredsFile,
		endpoint:  fmt.Sprintf("%s-documentai.googleapis.com:443", cfg.Location),
		name:      fmt.Sprintf("projects/%s/locations/%s/processors/%s", cfg.ProjectId, cfg.Location, cfg.ProcessorId),
	}
}

func (docAI *GoogleDocumentAI) Schema() string {
	return config.SCHEMA_DOCUMENT_AI
}

func (docAI *GoogleDocumentAI) Process(record database.Record, fileStore store.FileStore) error {
	ctx := context.Background()
	client, err := docAI.newDocumentProcessorClient(ctx)
	if err != nil {
		return err
	}
	req, err := docAI.newDocumentProcessorRequest(ctx, record, fileStore)
	if err != nil {
		return err
	}

	resp, err := client.ProcessDocument(ctx, req)
	if err != nil {
		return fmt.Errorf("failed GoogleDocumentAI ProcessDocument call: %w", err)
	}
	doc := resp.GetDocument()
	json, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed GoogleDocumentAI doc to json marshal: %w", err)
	}
	jsonPath := fmt.Sprintf("%s.json", record.Id)
	err = fileStore.Store(jsonPath, bytes.NewReader(json))
	if err != nil {
		return fmt.Errorf("failed GoogleDocumentAI file storage: %w", err)
	}

	return nil
}

func (docAI *GoogleDocumentAI) newDocumentProcessorClient(ctx context.Context) (*documentai.DocumentProcessorClient, error) {
	auth := option.WithCredentialsFile(docAI.credsFile)
	endpointOpt := option.WithEndpoint(docAI.endpoint)
	return documentai.NewDocumentProcessorClient(ctx, endpointOpt, auth)
}

func (docAI *GoogleDocumentAI) newDocumentProcessorRequest(ctx context.Context, record database.Record, fileStore store.FileStore) (*documentaipb.ProcessRequest, error) {
	f, err := fileStore.Get(record.Path)
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
