package dtos

type ForecastRequestParams struct {
	Lat  float64 `form:"lat" binding:"required"`
	Lon  float64 `form:"lon" binding:"required"`
	Date string  `form:"date" binding:"required"`
}

type ForecastResponse struct {
	TempC      float64 `json:"temp_c"`
	PrecipProb float64 `json:"precip_prob"`
	WindKph    float64 `json:"wind_kph"`
	Summary    string  `json:"summary"`
}
