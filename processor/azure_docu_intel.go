package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/store"
)

type DocuIntel struct {
	client     *http.Client
	endpoint   string
	key        string
	modelId    string
	apiVersion string
}

const (
	KEY_HEADER        = "Ocp-Apim-Subscription-Key"
	RESULT_ID_HEADER  = "Apim-Request-Id"
	MAX_FETCH_RETRIES = 5
)

func NewDocuIntel(cfg config.DocuIntelCfg) *DocuIntel {
	return &DocuIntel{
		client:     &http.Client{},
		key:        cfg.Key,
		endpoint:   cfg.Endpoint,
		modelId:    cfg.ModelId,
		apiVersion: cfg.ApiVersion,
	}
}

func (docInt *DocuIntel) Schema() string {
	return config.SCHEMA_DOC_INT
}

func (docInt *DocuIntel) Process(receipt database.Receipt, fs store.FileStore) error {
	req, err := docInt.newProcessRequest(receipt, fs)
	if err != nil {
		return err
	}

	res, err := docInt.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status not ok: %v", res.Status)
	}

	id := res.Header.Get(RESULT_ID_HEADER)
	if id == "" {
		return fmt.Errorf("could not retrieve id from response")
	}

	docInt.fetchResult(id, receipt, fs)

	return nil
}

func (docInt *DocuIntel) doRequest(req *http.Request) ([]byte, error) {
	res, err := docInt.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
	default:
		return nil, fmt.Errorf("status not ok: %v", res.Status)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return b, nil
}

func (docInt *DocuIntel) newProcessRequest(receipt database.Receipt, fileStore store.FileStore) (*http.Request, error) {
	url := fmt.Sprintf("%s/formrecognizer/documentModels/%s:analyze?api-version=%s", docInt.endpoint, docInt.modelId, docInt.apiVersion)

	type Payload struct {
		UrlSource string `json:"urlSource"`
	}

	sourceUrl, err := fileStore.GetURL(receipt.Path)
	if err != nil {
		return nil, err
	}

	payload := Payload{UrlSource: sourceUrl}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set(KEY_HEADER, docInt.key)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (docInt *DocuIntel) newResultRequest(resultId string) (*http.Request, error) {
	url := fmt.Sprintf("%s/formrecognizer/documentModels/%s/analyzeResults/%s?api-version=%s", docInt.endpoint, docInt.modelId, resultId, docInt.apiVersion)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(KEY_HEADER, docInt.key)
	return req, nil
}

func (docInt *DocuIntel) analyzeResults(resultId string) ([]byte, error) {
	req, err := docInt.newResultRequest(resultId)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	b, err := docInt.doRequest(req)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (docInt *DocuIntel) fetchResult(resultId string, receipt database.Receipt, fileStore store.FileStore) error {
	retries := MAX_FETCH_RETRIES
	var b []byte
	var err error
	for {
		if retries == 0 {
			return errors.New("failed DocInt fetchResult retries exceeded")
		}
		b, err = docInt.analyzeResults(resultId)
		if err != nil {
			return fmt.Errorf("failed DocInt analyzeResults: %w", err)
		}
		jMap := make(map[string]string, 0)
		json.Unmarshal(b, &jMap) // Ignore error. Only care about status field. Rest can fail.

		if jMap["status"] == "succeeded" {
			jsonPath := fmt.Sprintf("%s.json", receipt.Id)
			err = fileStore.Store(jsonPath, bytes.NewReader(b))
			if err != nil {
				return fmt.Errorf("failed DocInt store: %w", err)
			}
			break
		}
		retries--
		time.Sleep(time.Duration(MAX_FETCH_RETRIES-retries+1) * time.Second)
	}

	return nil
}
