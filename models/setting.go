package models

type SubmitKratosSettingRequest struct {
	Method string  `json:"method"`
	Traits *Traits `json:"traits"`
}
