package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mcp-server/internal/dtos"
	"net/http"
)

type CalendarService struct {
	Client *http.Client
}

func NewCalendarService(client *http.Client) *CalendarService {
	return &CalendarService{Client: client}
}

func (s *CalendarService) GetHolidays(params dtos.HolidayRequestParams) (dtos.HolidayResponse, error) {
	log.Printf("Service: Calling Nager.Date API for: %s, %d", params.CountryCode, params.Year)

	baseURL := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/%s",
		params.Year,
		params.CountryCode,
	)

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, errors.New("failed to create Nager.Date request")
	}
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.New("external Nager.Date API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("No holidays found for %s, %d", params.CountryCode, params.Year)
			return dtos.HolidayResponse{}, nil
		}
		return nil, fmt.Errorf("Nager.Date API returned status: %s", resp.Status)
	}

	var data dtos.HolidayResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("failed to decode Nager.Date response")
	}

	return data, nil
}
