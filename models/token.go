package models

type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
