package dtos

type HolidayRequestParams struct {
	CountryCode string `form:"country_code" binding:"required,len=2"`
	Year        int    `form:"year" binding:"required"`
}

type Holiday struct {
	Date      string `json:"date"`
	LocalName string `json:"localName"`
}

type HolidayResponse []Holiday
