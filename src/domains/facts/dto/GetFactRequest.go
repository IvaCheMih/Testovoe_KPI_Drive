package dto

type GetFactRequest struct {
	PeriodStart     string `json:"period_start" form:"period_start"`
	PeriodEnd       string `json:"period_end" form:"period_end"`
	PeriodKey       string `json:"period_key" form:"period_key"`
	IndicatorToMoId string `json:"indicator_to_mo_id" form:"indicator_to_mo_id"`
}

type GetFactResponse struct {
	Message string
}
