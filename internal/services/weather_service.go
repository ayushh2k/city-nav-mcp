package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mcp-server/internal/dtos"
	"net/http"
	"net/url"
)

type WeatherService struct {
	Client *http.Client
}

func NewWeatherService(client *http.Client) *WeatherService {
	return &WeatherService{Client: client}
}

type openMeteoResponse struct {
	Daily struct {
		Time       []string  `json:"time"`
		TempMax    []float64 `json:"temperature_2m_max"`
		PrecipProb []float64 `json:"precipitation_probability_mean"`
		WindSpeed  []float64 `json:"windspeed_10m_max"`
	} `json:"daily"`
}

func (s *WeatherService) GetForecast(params dtos.ForecastRequestParams) (*dtos.ForecastResponse, error) {
	log.Printf("Service: Calling external API for forecast: lat=%f, lon=%f\n", params.Lat, params.Lon)

	baseURL, _ := url.Parse("https://api.open-meteo.com/v1/forecast")
	query := baseURL.Query()
	query.Set("latitude", fmt.Sprintf("%f", params.Lat))
	query.Set("longitude", fmt.Sprintf("%f", params.Lon))
	query.Set("start_date", params.Date)
	query.Set("end_date", params.Date)
	query.Set("timezone", "auto")
	query.Set("daily", "temperature_2m_max,precipitation_probability_mean,windspeed_10m_max")
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		log.Printf("Error creating weather request: %v", err)
		return nil, errors.New("failed to create weather request")
	}

	req.Header.Set("User-Agent", "CityNavigator-MCP-Server (Go-http-client)")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error calling Open-Meteo: %v", err)
		return nil, errors.New("external weather API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Open-Meteo returned non-200 status: %s", resp.Status)
		return nil, fmt.Errorf("external API returned status: %s", resp.Status)
	}

	var weatherData openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		log.Printf("Error decoding weather JSON: %v", err)
		return nil, errors.New("failed to decode weather API response")
	}

	if len(weatherData.Daily.Time) == 0 ||
		len(weatherData.Daily.TempMax) == 0 ||
		len(weatherData.Daily.PrecipProb) == 0 ||
		len(weatherData.Daily.WindSpeed) == 0 {
		log.Printf("Incomplete data from Open-Meteo for date: %s", params.Date)
		return nil, errors.New("incomplete data from weather API")
	}

	temp := weatherData.Daily.TempMax[0]
	precip := weatherData.Daily.PrecipProb[0]
	summary := fmt.Sprintf("Max temp: %.1f°C, Precip: %.0f%%", temp, precip)

	return &dtos.ForecastResponse{
		TempC:      temp,
		PrecipProb: precip,
		WindKph:    weatherData.Daily.WindSpeed[0],
		Summary:    summary,
	}, nil
}
