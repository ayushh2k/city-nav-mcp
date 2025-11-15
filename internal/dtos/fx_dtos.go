package dtos

type FxRequestParams struct {
	From   string  `form:"from" binding:"required"`
	To     string  `form:"to" binding:"required"`
	Amount float64 `form:"amount" binding:"required"`
}

type FxResponse struct {
	Rate      float64 `json:"rate"`
	Converted float64 `json:"converted"`
}
