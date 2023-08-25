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

func (dt *DocuIntelTransform) ToCommon() (*CommonData, error) {
	obj := RawDocIntData{}
	err := json.Unmarshal(dt.Data, &obj)
	if err != nil {
		fmt.Println("Error parsing DocuIntel data:", err)
	}
	fmt.Printf("%+v\n", obj.AnalyzeResult.Documents)
	return &CommonData{}, nil
}

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
