package models

type Fact struct {
	PeriodStart         string `json:"period_start" form:"period_start"`
	PeriodEnd           string `json:"period_end" form:"period_end"`
	PeriodKey           string `json:"period_key" form:"period_key"`
	IndicatorToMoId     string `json:"indicator_to_mo_id" form:"indicator_to_mo_id"`
	IndicatorToMoFactId string `json:"indicator_to_mo_fact_id" form:"indicator_to_mo_fact_id"`
	Value               string `json:"value" form:"value"`
	FactTime            string `json:"fact_time" form:"fact_time"`
	IsPlan              string `json:"is_plan" form:"is_plan"`
	AuthUserId          string `json:"auth_user_id" form:"auth_user_id"`
	Comment             string `json:"comment" form:"comment"`
}
