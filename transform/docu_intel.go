package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/store"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

type DocuIntelTransform struct {
	Record database.Record
	Data   []byte
}

func NewDocuIntelTransform(data []byte, record database.Record) *DocuIntelTransform {
	return &DocuIntelTransform{
		Data: data,
	}
}

func (dt *DocuIntelTransform) ToCommon() (*Expense, error) {
	obj := RawDocIntData{}
	err := json.Unmarshal(dt.Data, &obj)
	if err != nil {
		fmt.Println("Error parsing DocuIntel data:", err)
	}
	expense := dt.mapFields(obj.AnalyzeResult.Documents[0].Fields)
	fmt.Printf("%+v\n", expense)

	data, err := json.Marshal(expense)
	if err != nil {
		log.Printf("failed to marshal Expense in Docu Intel Transform: %s", err)
	}

	store.File.Store("common-"+dt.Record.JSON, bytes.NewReader(data))

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
