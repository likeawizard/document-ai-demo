package transform

import (
	"encoding/json"
	"fmt"
	"time"
)

type DocuIntelTransform struct {
	Data []byte
}

func NewDocuIntelTransform(data []byte) *DocuIntelTransform {
	return &DocuIntelTransform{
		Data: data,
	}
}

func (dt *DocuIntelTransform) ToCommon() (*Expense, error) {
	obj := RawDocIntData{}
	err := json.Unmarshal(dt.Data, &obj)
	if err != nil {
		return nil, fmt.Errorf("error parsing DocuIntel data: %s", err)
	}
	expense := dt.mapFields(obj.AnalyzeResult.Documents[0].Fields)

	return &expense, nil
}

func (dt *DocuIntelTransform) mapFields(fields Fields) Expense {
	exp := Expense{}
	exp.Currency = fields.Currency.ValueString
	transTime, err := time.Parse(timeLayout, fmt.Sprintf("%s %s", fields.TransactionDate.ValueDate, fields.TransactionTime.ValueTime))
	if err == nil {
		exp.Date = transTime
	}
	exp.Total = fields.Total.ValueNumber
	exp.Merchant.MerchantName = fields.MerchantName.ValueString
	exp.Merchant.MerchantAddress = fields.MerchantAddress.ValueAddress.String()
	exp.Merchant.MerchantPhone = fields.MerchantPhoneNumber.ValuePhoneNumber

	return exp
}
