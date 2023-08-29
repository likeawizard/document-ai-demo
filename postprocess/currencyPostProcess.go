package postprocess

import (
	"fmt"
	"time"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/transform"
)

const (
	TARGET_CURRENCY = "EUR"
)

type CurrencyService interface {
	GetConversionRate(*CurrencyPostProcess) error
}

type CurrencyPostProcess struct {
	rate     float64
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

func (pp *CurrencyPostProcess) GetFields(exp transform.Expense) error {
	switch {
	case exp.Currency == TARGET_CURRENCY:
		return fmt.Errorf("nothing to convert")
	case len(exp.Currency) != 3:
		return fmt.Errorf("currency not in the three letter ISO format: '%s'", exp.Currency)
	case exp.Total == 0:
		return fmt.Errorf("nothing to convert total is zero")
	}

	if exp.Date.IsZero() {
		// TODO historic rates are available only for dates until last midnight. Should be handled in API call - set past date there or query live currency rate
		pp.date = time.Now().Add(-1 * 24 * time.Hour)
	} else {
		pp.date = exp.Date
	}

	pp.currency = exp.Currency

	pp.fields = FieldMap{
		"total": fmt.Sprintf("%f", exp.Total),
		"tax":   fmt.Sprintf("%f", exp.Tax),
	}

	return nil
}

func (pp *CurrencyPostProcess) Apply(exp *transform.Expense) {
	for k := range pp.fields {
		switch k {
		case "tax":
			exp.Tax *= pp.rate
		case "total":
			exp.Total *= pp.rate
		}
	}

	exp.Currency = TARGET_CURRENCY
}
