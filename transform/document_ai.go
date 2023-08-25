package transform

import (
	"encoding/json"
	"fmt"
)

type DocumentAiTransform struct {
	Data []byte
}

func NewDocumentAiTransform(data []byte) *DocumentAiTransform {
	return &DocumentAiTransform{
		Data: data,
	}
}

func (dt *DocumentAiTransform) ToCommon() (*CommonData, error) {
	obj := RawDocumentAiData{}
	err := json.Unmarshal(dt.Data, &obj)
	if err != nil {
		fmt.Println("Error parsing DocumentAi data:", err)
	}
	fmt.Printf("%+v\n", obj.Entities)
	return &CommonData{}, nil
}

type RawDocumentAiData struct {
	Text     string   `json:"text"`
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Type            string          `json:"type"`
	MentionText     string          `json:"mention_text,omitempty"`
	Confidence      float64         `json:"confidence"`
	ID              string          `json:"id"`
	NormalizedValue NormalizedValue `json:"normalized_value,omitempty"`
	Properties      []Property      `json:"properties,omitempty"`
}

type Property struct {
	Type        string  `json:"type"`
	MentionText string  `json:"mention_text"`
	Confidence  float64 `json:"confidence"`
	ID          string  `json:"id"`
}

type DateValue struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type DatetimeValue struct {
	Hours      int         `json:"hours"`
	Minutes    int         `json:"minutes"`
	TimeOffset interface{} `json:"TimeOffset"`
}

type StructuredValue struct {
	DateValue     DateValue     `json:"DateValue"`
	DatetimeValue DatetimeValue `json:"DatetimeValue"`
}

type NormalizedValue struct {
	StructuredValue StructuredValue `json:"StructuredValue"`
	Text            string          `json:"text"`
}
