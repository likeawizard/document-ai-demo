package postprocess

import (
	"fmt"

	"github.com/likeawizard/document-ai-demo/transform"
	"golang.org/x/text/language"
)

type TranslationPostProcess struct {
	lang         language.Tag
	translations FieldMap
	fields       FieldMap
}

func (pp *TranslationPostProcess) GetFields(exp transform.Expense) error {
	// TODO: detect document language from processor / transform stage and only apply translation post process if not english
	detectLang := "any"
	if detectLang == "en" {
		return fmt.Errorf("nothing to translate document langauge: '%s'", detectLang)
	}
	pp.lang = language.English

	pp.fields = FieldMap{
		"merchantAddr": exp.Merchant.MerchantAddress,
		"merchant":     exp.Merchant.MerchantName,
		"merchantStr":  exp.Merchant.StringVal,
	}

	return nil
}

func (pp *TranslationPostProcess) Apply(exp *transform.Expense) {
	for k, v := range pp.fields {
		switch k {
		case "merchantAddr":
			exp.Merchant.MerchantAddress = v
		case "merchant":
			exp.Merchant.MerchantName = v
		case "merchantStr":
			exp.Merchant.StringVal = v
		}
	}

	exp.Currency = TARGET_CURRENCY
}
