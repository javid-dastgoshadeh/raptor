package models

type SubmitKratosRegisterRequest struct {
	Method   string  `json:"method"`
	Password string  `json:"password"`
	Traits   *Traits `json:"traits"`
}

type AppPreRegisterResponse struct {
	M        string  `json:"method"`
	Password string  `json:"password"`
	Traits   *Traits `json:"traits"`
}
