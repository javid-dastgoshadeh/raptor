package account

import (
	"context"
	"raptor/models"
)

// KratosRepository ...
type KratosRepository interface {
	CreateSettingFlow(ctx context.Context, claims map[string]interface{}) (string, error)
	SubmitSettingFlow(ctx context.Context, identityTraits *models.SubmitKratosSettingRequest, flow string, claims map[string]interface{}) (interface{}, error)
	CheckSession(ctx context.Context, claims map[string]interface{}) (interface{}, error)
	DisableSession(ctx context.Context, claims map[string]interface{}) error
	InactiveIdentity(ctx context.Context, claims map[string]interface{}) error
	GetIdentity(ctx context.Context, claims map[string]interface{}) (interface{}, error)
	UpdateIdentity(ctx context.Context, identityID string, identity interface{}) (interface{}, error)
}

// IdentificationRepository ...
type IdentificationRepository interface {
	CheckIdentification(ctx context.Context, jwt string) (bool, error)
}
