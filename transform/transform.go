package transform

import (
	"fmt"
	"strconv"
	"time"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
)

type DataTransform interface {
	ToCommon() (*Expense, error)
}

func NewDataTransform(schema string, data []byte, record database.Record) (DataTransform, error) {
	switch schema {
	case config.SCHEMA_DOCUMENT_AI:
		return NewDocumentAiTransform(data, record), nil
	case config.SCHEMA_DOC_INT:
		return NewDocuIntelTransform(data, record), nil
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
