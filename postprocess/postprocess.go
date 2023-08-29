package postprocess

import "github.com/likeawizard/document-ai-demo/transform"

type FieldMap map[string]string

type PostProcessor interface {
	GetFields(transform.Expense)
	PostProcess() error
	Apply(*transform.Expense)
}
