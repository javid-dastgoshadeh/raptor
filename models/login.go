package models

type SubmitKratosLoginRequest struct {
	Method     string `json:"method"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type CodeRequest struct {
	EmailMobile string `json:"email_mobile" validate:"required"`
}

type VerifyCodeRequest struct {
	EmailMobile string `json:"email_mobile" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Device      string `json:"device,omitempty"`
}

type AfterRegisterRequest struct {
	EmailMobile string `json:"email_mobile" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
}
