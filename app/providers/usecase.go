package providers

import (
	"context"
	"raptor/models"
)

// OtpUsecase ...
type OtpUsecase interface {
	SendOtpCode(ctx context.Context, emailMobile string) error
	VerifyOtpCode(ctx context.Context, emailMobile string, code string, device string) (*models.TokenResponse, error)
	AfterRegisterStep(ctx context.Context, emailMobile string, code string, displayName string) (*models.TokenResponse, error)
}

// SocialUsecase ...
type SocialUsecase interface {
	VerifySocialCode(ctx context.Context, provider string, data map[string]interface{}) (*models.TokenResponse, error)
	ExchangeSocialCode(ctx context.Context, provider string, data map[string]interface{}) (*models.TokenResponse, error)
}

// TokenUsecase ...
type TokenUsecase interface {
	RefreshToken(ctx context.Context, jwtClaims map[string]interface{}) (*models.TokenResponse, error)
}
