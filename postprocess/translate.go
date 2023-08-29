package postprocess

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type TranslationService interface {
	Translate(*TranslationPostProcess) error
}

type GoogleTranslationService struct {
	authFile string
}

func NewTranslationSerivce() (TranslationService, error) {
	// TODO: hardcoded creds
	return &GoogleTranslationService{authFile: "document-ai-creds.json"}, nil
}

func (ts *GoogleTranslationService) Translate(tpp *TranslationPostProcess) error {
	ctx := context.Background()

	var keys, vals []string
	for k, v := range tpp.fields {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	// TODO: hardcoded credentials file.
	auth := option.WithCredentialsFile(ts.authFile)
	client, err := translate.NewClient(ctx, auth)
	if err != nil {
		return fmt.Errorf("failed to translation initialize client: %v", err)
	}
	lang := language.English
	translation, err := client.Translate(ctx, vals, lang, nil)
	if err != nil {
		return fmt.Errorf("failed to translate: %v", err)
	}
	translated := make(FieldMap)

	for i := range keys {
		translated[keys[i]] = translation[i].Text
	}
	tpp.translations = translated

	return nil
}
