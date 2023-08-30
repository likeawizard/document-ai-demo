package transform

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	dateType            = "receipt_date"
	timeType            = "purchase_time"
	supplierNameType    = "supplier_name"
	supplierAddressType = "supplier_address"
	totalType           = "total_amount"
	currecnyType        = "currency"
)

type DocumentAiTransform struct {
	Data []byte
}

func NewDocumentAiTransform(data []byte) *DocumentAiTransform {
	return &DocumentAiTransform{
		Data: data,
	}
}

func (dt *DocumentAiTransform) ToCommon() (*Expense, error) {
	obj := RawDocumentAiData{}
	err := json.Unmarshal(dt.Data, &obj)
	if err != nil {
		return nil, fmt.Errorf("error parsing DocumentAi data: %s", err)
	}
	expense := dt.mapFields(obj.Entities)

	return &expense, nil
}

func (dt *DocumentAiTransform) mapFields(entities []Entity) Expense {
	expense := Expense{}
	var dateStr, timeStr string
	for _, entity := range entities {
		switch entity.Type {
		case dateType:
			date := entity.NormalizedValue.StructuredValue.DateValue
			if date.Year != 0 && date.Month != 0 && date.Day != 0 {
				dateStr = fmt.Sprintf("%d-%02d-%02d", date.Year, date.Month, date.Day)
			}
		case timeType:
			t := entity.NormalizedValue.StructuredValue.DatetimeValue
			if t.Hours != 0 && t.Minutes != 0 {
				timeStr = fmt.Sprintf("%02d:%02d:00", t.Hours, t.Minutes)
			}
		case supplierNameType:
			expense.Merchant.MerchantName = entity.MentionText
		case supplierAddressType:
			expense.Merchant.MerchantAddress = entity.NormalizedValue.Text
		case totalType:
			expense.Total = moneyParser(entity.MentionText)
		case currecnyType:
			expense.Currency = entity.NormalizedValue.Text
		}
	}

	if dateStr != "" {
		if timeStr != "" {
			dateStr = fmt.Sprintf("%s %s", dateStr, timeStr)
		}
		datetime, err := time.Parse(timeLayout, dateStr)
		if err == nil {
			expense.Date = datetime
		}
	}

	return expense
}
