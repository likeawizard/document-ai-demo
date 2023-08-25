package transform

import (
	"fmt"

	"github.com/likeawizard/document-ai-demo/config"
)

type DataTransform interface {
	ToCommon() (*CommonData, error)
}

func NewDataTransform(schema string, data []byte) (DataTransform, error) {
	switch schema {
	case config.SCHEMA_DOCUMENT_AI:
		return NewDocumentAiTransform(data), nil
	case config.SCHEMA_DOC_INT:
		return NewDocuIntelTransform(data), nil
	default:
		return nil, fmt.Errorf("unsupported document schema for transform: '%s'", schema)
	}
}

type CommonData struct {
	// TODO - implement generic receipt data structure that different schemas can transform to
}
