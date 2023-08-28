package transform

import (
	"fmt"
	"time"
)

type RawDocIntData struct {
	Status              string        `json:"status"`
	SreatedDateTime     time.Time     `json:"createdDateTime"`
	SastUpdatedDateTime time.Time     `json:"lastUpdatedDateTime"`
	AnalyzeResult       AnalyzeResult `json:"analyzeResult"`
}

type AnalyzeResult struct {
	APIVersion string     `json:"apiVersion"`
	ModelID    string     `json:"modelId"`
	Content    string     `json:"content"`
	Documents  []Document `json:"documents"`
}

type Document struct {
	DocType    string  `json:"docType"`
	Fields     Fields  `json:"fields"`
	Confidence float64 `json:"confidence"`
}

type Fields struct {
	Currency            Currency            `json:"Currency,omitempty"`
	MerchantName        MerchantName        `json:"MerchantName,omitempty"`
	MerchantAddress     MerchantAddress     `json:"MerchantAddress,omitempty"`
	MerchantPhoneNumber MerchantPhoneNumber `json:"MerchantPhoneNumber,omitempty"`
	TransactionDate     TransactionDate     `json:"TransactionDate,omitempty"`
	TransactionTime     TransactionTime     `json:"TransactionTime,omitempty"`
	Subtotal            Subtotal            `json:"Subtotal,omitempty"`
	Total               Total               `json:"Total,omitempty"`
}

type Currency struct {
	Type        string  `json:"type"`
	ValueString string  `json:"valueString"`
	Content     string  `json:"content"`
	Confidence  float64 `json:"confidence"`
}

type MerchantName struct {
	Type        string  `json:"type"`
	ValueString string  `json:"valueString"`
	Content     string  `json:"content"`
	Confidence  float64 `json:"confidence"`
}

type Total struct {
	Type        string  `json:"type"`
	ValueNumber float64 `json:"valueNumber"`
	Content     string  `json:"content"`
	Confidence  float64 `json:"confidence"`
}

type Subtotal struct {
	Type        string  `json:"type"`
	ValueNumber float64 `json:"valueNumber"`
	Content     string  `json:"content"`
	Confidence  float64 `json:"confidence"`
}

type MerchantAddress struct {
	Type         string       `json:"type"`
	Content      string       `json:"content"`
	Confidence   float64      `json:"confidence"`
	ValueAddress ValueAddress `json:"valueAddress"`
}

type ValueAddress struct {
	HouseNumber   string `json:"houseNumber"`
	Road          string `json:"road"`
	PostalCode    string `json:"postalCode"`
	CountryRegion string `json:"countryRegion"`
	StreetAddress string `json:"streetAddress"`
	Unit          string `json:"unit"`
	CityDistrict  string `json:"cityDistrict"`
}

type MerchantPhoneNumber struct {
	Type             string  `json:"type"`
	ValuePhoneNumber string  `json:"valuePhoneNumber"`
	Content          string  `json:"content"`
	Confidence       float64 `json:"confidence"`
}

type TransactionDate struct {
	Type       string  `json:"type"`
	ValueDate  string  `json:"valueDate"`
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"`
}

type TransactionTime struct {
	Type       string  `json:"type"`
	ValueTime  string  `json:"valueTime"`
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"`
}

// TODO: improve address formatting when fields are zero valued
func (val ValueAddress) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", val.StreetAddress, val.CityDistrict, val.CountryRegion, val.PostalCode)
}
