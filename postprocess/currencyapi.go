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

func (cs *CurrencyApiService) GetConversionRate(from, to string, date time.Time) (float64, error) {
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
	params.Add("date", date.Format("2006-01-02"))
	params.Add("base_currency", from)
	params.Add("currencies", to)
	params.Add("apikey", cs.authKey)

	u, err := url.ParseRequestURI(cs.endpoint)
	if err != nil {
		return 1, err
	}

	u.Path = cs.getHistoricalPath()
	u.RawQuery = params.Encode()
	fmt.Println(u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 1, err
	}
	req.Header.Add("apikey", cs.authKey)

	resp, err := cs.client.Do(req)
	if err != nil {
		return 1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return 1, err
	}

	var currencyData response
	err = json.Unmarshal(body, &currencyData)
	if err != nil {
		return 1, err
	}

	fmt.Println(currencyData)

	conversionData, ok := currencyData.Data[to]
	if !ok {
		return 1, fmt.Errorf("requested currency conversion missing '%s'", to)
	}

	return conversionData.Value, nil
}
