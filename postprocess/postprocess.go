package postprocess

import (
	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/store"
	"github.com/likeawizard/document-ai-demo/transform"
)

type FieldMap map[string]string

type PostProcessor interface {
	GetFields(transform.Expense) error
	Apply(*transform.Expense)
}

type PostProcessService struct {
	CurrencyService    CurrencyService
	TranslationService TranslationService
	FileStore          store.FileStore
}

func NewPostProcessService(cfg config.Config) (*PostProcessService, error) {
	pps := PostProcessService{}

	cs, err := NewCurrencyService(cfg.Currency)
	if err != nil {
		return nil, err
	}
	pps.CurrencyService = cs

	ts, err := NewTranslationSerivce()
	if err != nil {
		return nil, err
	}
	pps.TranslationService = ts

	fs, err := store.NewFileStore(cfg.Store)
	if err != nil {
		return nil, err
	}
	pps.FileStore = fs

	return &pps, nil
}

func (pps *PostProcessService) GetCurrencyPostProcess(exp transform.Expense) (*CurrencyPostProcess, error) {
	cpp := CurrencyPostProcess{}
	err := cpp.GetFields(exp)
	if err != nil {
		return nil, err
	}

	return &cpp, nil
}

func (pps *PostProcessService) GetTranslationPostProcess(exp transform.Expense) (*TranslationPostProcess, error) {
	tpp := TranslationPostProcess{}
	err := tpp.GetFields(exp)
	if err != nil {
		return nil, err
	}

	return &tpp, nil
}
