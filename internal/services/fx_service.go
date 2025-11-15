package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mcp-server/internal/dtos"
	"net/http"
	"net/url"
	"strings"
)

type FxService struct {
	Client *http.Client
	APIKey string
}

func NewFxService(client *http.Client, apiKey string) *FxService {
	return &FxService{
		Client: client,
		APIKey: apiKey,
	}
}

type exchangerateLatestResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func (s *FxService) ConvertCurrency(params dtos.FxRequestParams) (*dtos.FxResponse, error) {
	log.Printf("Service: Calling exchangeratesapi.io /latest for %s,%s", params.From, params.To)

	baseURL, _ := url.Parse("https://api.exchangeratesapi.io/v1/latest")
	query := baseURL.Query()
	query.Set("access_key", s.APIKey)
	query.Set("symbols", fmt.Sprintf("%s,%s", params.From, params.To))
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, errors.New("failed to create fx request")
	}
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.New("external fx API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exchangeratesapi.io API returned status: %s", resp.Status)
	}

	var data exchangerateLatestResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("failed to decode fx response")
	}

	rateFrom, ok := data.Rates[strings.ToUpper(params.From)]
	if !ok {
		return nil, fmt.Errorf("currency %s not found", params.From)
	}

	rateTo, ok := data.Rates[strings.ToUpper(params.To)]
	if !ok {
		return nil, fmt.Errorf("currency %s not found", params.To)
	}

	finalRate := rateTo / rateFrom
	finalConverted := params.Amount * finalRate

	return &dtos.FxResponse{
		Rate:      finalRate,
		Converted: finalConverted,
	}, nil
}
