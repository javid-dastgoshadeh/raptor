package providers

import (
	"context"
	"raptor/models"
)

// KratosRepository ...
type KratosRepository interface {
	CreateKratosLoginFLow(ctx context.Context) (flow string, err error)
	SubmitKratosLoginFlow(ctx context.Context, flow string, req *models.SubmitKratosLoginRequest) (interface{}, error)
	CreateKratosRegisterFLow(ctx context.Context) (flow string, err error)
	SubmitKratosRegisterFLow(ctx context.Context, flow string, req *models.SubmitKratosRegisterRequest) (interface{}, error)
	CheckIdentityExistence(ctx context.Context, identity string, identityType string) (bool, error)
	UpdateIdentity(ctx context.Context, identityID string, identity interface{}) (interface{}, error)
	GetIdentity(ctx context.Context, claims map[string]interface{}) (interface{}, error)
	ActiveIdentity(ctx context.Context, jwtClaims map[string]interface{}, identityFormat string) (interface{}, error)
}

// OtpRepository ...
type OtpRepository interface {
	SendCode(ctx context.Context, identity string, code string, senderType *models.MessageSender) error
	VerifyCode(ctx context.Context, identity string, code string) error
}

// SocialRepository ...
type SocialRepository interface {
	VerifyGoogle(ctx context.Context, data map[string]interface{}) (*models.Traits, error)
	ExchangeGoogleOauth2Code(ctx context.Context, data map[string]interface{}) (*models.Traits, error)
	VerifyApple(ctx context.Context, data map[string]interface{}) (*models.Traits, error)
	ExchangeAppleOauth2Code(ctx context.Context, data map[string]interface{}) (*models.Traits, error)
}
