package account

import (
	"context"
	"raptor/models"
)

// Usecase ...
type Usecase interface {
	Logout(ctx context.Context, jwtClaims map[string]interface{}) error
	Profile(ctx context.Context, jwtClaims map[string]interface{}) (interface{}, error)
	UpdateProfile(ctx context.Context, jwtClaims map[string]interface{}, request *models.UpdateRequest) (interface{}, error)
	UpdateIdentifier(ctx context.Context, jwt string, request map[string]interface{}) (interface{}, error)
	VerifyUpdateIdentifier(ctx context.Context, jwtClaims map[string]interface{}, identityTraits map[string]interface{}) (interface{}, error)
	InactivateAccount(ctx context.Context, jwtClaims map[string]interface{}) error
}
