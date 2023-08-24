package expensebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
)

type DocuIntel struct {
	client     *http.Client
	endpoint   string
	key        string
	modelId    string
	apiVersion string
}

const (
	KEY_HEADER       = "Ocp-Apim-Subscription-Key"
	RESULT_ID_HEADER = "Apim-Request-Id"
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

func (docInt *DocuIntel) Process(record database.Record) error {
	req, err := docInt.newProcessRequest(record)
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
	fmt.Printf("result id: %v\n", id)
	// TODO

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

func (docInt *DocuIntel) newProcessRequest(record database.Record) (*http.Request, error) {
	url := fmt.Sprintf("%s/formrecognizer/documentModels/%s:analyze?api-version=%s", docInt.endpoint, docInt.modelId, docInt.apiVersion)

	type Payload struct {
		UrlSource string `json:"urlSource"`
	}
	//TODO get actual link from store
	payload := Payload{UrlSource: "https://storage.googleapis.com/reciept-store/receipt5.png"}
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

func (docInt *DocuIntel) analyzeResults(resultId string) error {
	req, err := docInt.newResultRequest(resultId)
	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
	}

	b, err := docInt.doRequest(req)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	// TODO
	return nil
}
