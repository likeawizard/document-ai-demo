package postprocess

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/likeawizard/document-ai-demo/config"
)

type CurrencyApiService struct {
	client   *http.Client
	endpoint string
	authKey  string
}

func NewCurrencyApiService(cfg config.CurrencyCfg) *CurrencyApiService {
	cs := CurrencyApiService{
		client:   &http.Client{},
		endpoint: cfg.Endpoint,
		authKey:  cfg.AuthKey,
	}

	return &cs
}

func (cs *CurrencyApiService) getHistoricalPath() string {
	return "/v3/historical"
}

func (cs *CurrencyApiService) GetConversionRate(cpp *CurrencyPostProcess) error {
	type response struct {
		Meta struct {
			LastUpdatedAt time.Time `json:"last_updated_at"`
		} `json:"meta"`
		Data map[string]struct {
			Code  string  `json:"code"`
			Value float64 `json:"value"`
		} `json:"data"`
	}

	params := url.Values{}
	params.Add("date", cpp.date.Format("2006-01-02"))
	params.Add("base_currency", cpp.currency)
	params.Add("currencies", TARGET_CURRENCY)
	params.Add("apikey", cs.authKey)

	u, err := url.ParseRequestURI(cs.endpoint)
	if err != nil {
		return err
	}

	u.Path = cs.getHistoricalPath()
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("apikey", cs.authKey)

	resp, err := cs.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var currencyData response
	err = json.Unmarshal(body, &currencyData)
	if err != nil {
		return err
	}

	conversionData, ok := currencyData.Data[TARGET_CURRENCY]
	if !ok {
		return fmt.Errorf("requested currency conversion missing '%s'", TARGET_CURRENCY)
	}

	cpp.rate = conversionData.Value

	return nil
}
