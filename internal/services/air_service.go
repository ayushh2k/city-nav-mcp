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

type AirService struct {
	Client *http.Client
	APIKey string
}

func NewAirService(client *http.Client, apiKey string) *AirService {
	return &AirService{
		Client: client,
		APIKey: apiKey,
	}
}

type locationSensor struct {
	ID        int `json:"id"`
	Parameter struct {
		Name string `json:"name"`
	} `json:"parameter"`
}

type locationResponse struct {
	Results []struct {
		Sensors []locationSensor `json:"sensors"`
	} `json:"results"`
}

type measurementResponse struct {
	Results []struct {
		Value float64 `json:"value"`
	} `json:"results"`
}

func (s *AirService) GetAQI(params dtos.AirQualityRequestParams) (*dtos.AirQualityResponse, error) {
	log.Printf("Service: Calling OpenAQ v3 API: lat=%f, lon=%f", params.Lat, params.Lon)

	sensorMap, err := s.findClosestSensors(params.Lat, params.Lon)
	if err != nil {
		return nil, fmt.Errorf("could not find sensors: %v", err)
	}

	if len(sensorMap) == 0 {
		log.Printf("No sensors found at location: lat=%f, lon=%f", params.Lat, params.Lon)
		return nil, errors.New("no sensors found at that location")
	}

	results := make(map[string]float64)

	if sensorID, ok := sensorMap["pm25"]; ok {
		val, err := s.fetchSensorMeasurement(sensorID)
		if err != nil {
			log.Printf("Error fetching pm25: %v", err)
		}
		results["pm25"] = val
	}

	if sensorID, ok := sensorMap["pm10"]; ok {
		val, err := s.fetchSensorMeasurement(sensorID)
		if err != nil {
			log.Printf("Error fetching pm10: %v", err)
		}
		results["pm10"] = val
	}

	if sensorID, ok := sensorMap["no2"]; ok {
		val, err := s.fetchSensorMeasurement(sensorID)
		if err != nil {
			log.Printf("Error fetching no2: %v", err)
		}
		results["no2"] = val
	}

	if sensorID, ok := sensorMap["o3"]; ok {
		val, err := s.fetchSensorMeasurement(sensorID)
		if err != nil {
			log.Printf("Error fetching o3: %v", err)
		}
		results["o3"] = val
	}

	response := dtos.AirQualityResponse{
		PM25: results["pm25"],
		PM10: results["pm10"],
		NO2:  results["no2"],
		O3:   results["o3"],
	}
	response.Category = getAQICategory(response.PM25)

	return &response, nil
}

func (s *AirService) findClosestSensors(lat, lon float64) (map[string]int, error) {
	baseURL, _ := url.Parse("https://api.openaq.org/v3/locations")
	query := url.Values{}
	query.Set("coordinates", fmt.Sprintf("%f,%f", lat, lon))
	query.Set("radius", "50000")
	query.Set("limit", "1")
	baseURL.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", baseURL.String(), nil)
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API (locations) returned %s", resp.Status)
	}

	var data locationResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		return nil, errors.New("no locations found")
	}

	sensorMap := make(map[string]int)
	for _, sensor := range data.Results[0].Sensors {
		sensorMap[sensor.Parameter.Name] = sensor.ID
	}

	return sensorMap, nil
}

func (s *AirService) fetchSensorMeasurement(sensorID int) (float64, error) {
	baseURL, _ := url.Parse(fmt.Sprintf("https://api.openaq.org/v3/sensors/%d/measurements", sensorID))
	query := url.Values{}
	query.Set("order_by", "datetime")
	query.Set("sort", "desc")
	query.Set("limit", "1")
	baseURL.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", baseURL.String(), nil)
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server")

	resp, err := s.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API (sensors) returned %s", resp.Status)
	}

	var data measurementResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if len(data.Results) == 0 {
		return 0, errors.New("no measurements for this sensor")
	}

	return data.Results[0].Value, nil
}

func getAQICategory(pm25 float64) string {
	switch {
	case pm25 <= 12.0:
		return "Good"
	case pm25 <= 35.4:
		return "Moderate"
	case pm25 <= 55.4:
		return "Unhealthy for Sensitive Groups"
	case pm25 <= 75.0:
		return "Unhealthy"
	case pm25 <= 150.4:
		return "Very Unhealthy"
	default:
		return "Hazardous"
	}
}
