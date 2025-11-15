package dtos

type AirQualityRequestParams struct {
	Lat  float64 `form:"lat" binding:"required"`
	Lon  float64 `form:"lon" binding:"required"`
	Date string  `form:"date"`
}

type AirQualityResponse struct {
	PM25     float64 `json:"pm25"`
	PM10     float64 `json:"pm10"`
	NO2      float64 `json:"no2"`
	O3       float64 `json:"o3"`
	Category string  `json:"category"`
}
