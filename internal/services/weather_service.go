package services

import (
	"fmt"
	"mcp-server/internal/dtos"
	"net/http"
)

type WeatherService struct {
	Client *http.Client
}

func NewWeatherService(client *http.Client) *WeatherService {
	return &WeatherService{Client: client}
}

func (s *WeatherService) GetForecast(params dtos.ForecastRequestParams) (*dtos.ForecastResponse, error) {
	fmt.Printf("Service: Received forecast request (placeholder): %+v\n", params)

	return &dtos.ForecastResponse{
		TempC:      15.5,
		PrecipProb: 10.0,
		WindKph:    5.2,
		Summary:    "Partly cloudy (placeholder)",
	}, nil
}
