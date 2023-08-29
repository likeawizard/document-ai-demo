package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/store"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

type DataTransform interface {
	ToCommon() (*Expense, error)
}

type DataTransformService struct {
	FileStore store.FileStore
}

func NewDataTransformService(cfg config.Config) (*DataTransformService, error) {
	dts := DataTransformService{}
	fs, err := store.NewFileStore(cfg.Store)
	if err != nil {
		return nil, err
	}
	dts.FileStore = fs
	return &dts, nil
}

func (dts *DataTransformService) Transform(record database.Record, schema string) error {
	r, err := dts.FileStore.Get(record.JSON)
	if err != nil {
		return err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	dt, err := NewDataTransform(schema, data)
	if err != nil {
		return err
	}

	expense, err := dt.ToCommon()
	if err != nil {
		return err
	}

	data, err = json.Marshal(expense)
	if err != nil {
		return err
	}

	err = dts.FileStore.Store(record.JSON, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}

func NewDataTransform(schema string, data []byte) (DataTransform, error) {
	switch schema {
	case config.SCHEMA_DOCUMENT_AI:
		return NewDocumentAiTransform(data), nil
	case config.SCHEMA_DOC_INT:
		return NewDocuIntelTransform(data), nil
	default:
		return nil, fmt.Errorf("unsupported document schema for transform: '%s'", schema)
	}
}

type Expense struct {
	Date     time.Time `json:"date"`
	Currency string    `json:"currency"`
	Total    float64   `json:"total"`
	Tax      float64   `json:"tax"`
	Merchant Merchant  `json:"merchant"`
}

type Merchant struct {
	StringVal            string `json:"string_val"`
	MerchantName         string `json:"name"`
	MerchantRegistration string `json:"reg_no"`
	MerchantAddress      string `json:"address"`
	MerchantPhone        string `json:"phone"`
}

// TODO: discover formatting. check 3rd to last symbol to determine decimal separator
// remove everything else, replace discovered decimal with period and convert
func moneyParser(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return f
	}

	return f
}
