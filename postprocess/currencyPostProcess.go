package postprocess

import (
	"fmt"
	"strconv"
	"time"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/transform"
)

const (
	convertTo = "EUR"
)

type CurrencyService interface {
	GetConversionRate(from, to string, date time.Time) (float64, error)
}

type CurrencyPostProcess struct {
	date     time.Time
	currency string
	fields   FieldMap
}

func NewCurrencyService(cfg config.CurrencyCfg) (CurrencyService, error) {
	switch cfg.Service {
	case config.CUR_CURR_API:
		return NewCurrencyApiService(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported currency service: '%s'", cfg.Service)
	}
}

func (pp *CurrencyPostProcess) GetFields(exp transform.Expense) {
	if exp.Date.IsZero() {
		pp.date = time.Now()
	} else {
		pp.date = exp.Date
	}

	pp.currency = exp.Currency

	pp.fields = FieldMap{
		"currency": exp.Currency,
		"total":    fmt.Sprintf("%f", exp.Total),
		"tax":      fmt.Sprintf("%f", exp.Tax),
	}
}

func (pp *CurrencyPostProcess) Apply(exp *transform.Expense) {
	for k, v := range pp.fields {
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		switch k {
		case "tax":
			exp.Tax = floatVal
		case "total":
			exp.Total = floatVal
		}
	}

	exp.Currency = convertTo
}

func (pp *CurrencyPostProcess) PostProcess() error {
	return nil
}
