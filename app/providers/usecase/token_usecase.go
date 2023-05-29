package usecase

import (
	"context"
	"fmt"
	"raptor/app/providers"

	env "raptor/config"
	"raptor/models"
	"raptor/pkg/helpers"
)

// socialUsecase ...
type tokenUsecase struct {
	KratosRepo providers.KratosRepository
}

// NewToken ...
func NewToken(kratosRepo providers.KratosRepository) *tokenUsecase {
	return &tokenUsecase{
		KratosRepo: kratosRepo,
	}
}

// RefreshToken ...
func (u *tokenUsecase) RefreshToken(ctx context.Context, jwtClaims map[string]interface{}) (*models.TokenResponse, error) {
	var emailMobile string

	getIdentity, err := u.KratosRepo.GetIdentity(ctx, jwtClaims)
	if err != nil {
		return nil, err
	}
	//export traits from getIdentity response
	exportTraits, err := helpers.GetTraitsFromKratosResponse(getIdentity)
	if exportTraits.PhoneNumber != nil {
		emailMobile = fmt.Sprintf("%v", exportTraits.PhoneNumber)
	}

	if exportTraits.Email != nil {
		emailMobile = fmt.Sprintf("%v", exportTraits.PhoneNumber)
	}

	//if traitsInfo["phone_number"] != nil {
	//	emailMobile = fmt.Sprintf("%v", traitsInfo["phone_number"])
	//}
	//if traitsInfo["email"] != nil {
	//	emailMobile = fmt.Sprintf("%v", traitsInfo["email"])
	//}

	flow, err := u.KratosRepo.CreateKratosLoginFLow(ctx)

	if err != nil {
		return nil, err
	}

	submitLoginReq := &models.SubmitKratosLoginRequest{
		Method:     "password",
		Identifier: emailMobile,
		Password:   env.GetString("identity.password"),
	}
	res, err := u.KratosRepo.SubmitKratosLoginFlow(ctx, flow, submitLoginReq)

	claims, err := helpers.GenerateAccessTokenClaims(res)
	if err != nil {
		return nil, err
	}

	//load private key for access_token
	privateKey := models.PrivateKey
	accessToken, err := helpers.GenerateJwtToken(privateKey, claims)
	if err != nil {
		return nil, err
	}

	//generate refresh_token

	//parse response to get desire claims for jwt token
	refreshTokenClaims, err := helpers.GenerateRefreshTokenClaims(res)
	if err != nil {
		return nil, err
	}
	//load private key for refresh_token
	refreshTokenPrivateKey := models.RefreshTokenPrivateKey

	refreshToken, err := helpers.GenerateJwtToken(refreshTokenPrivateKey, refreshTokenClaims)
	if err != nil {
		return nil, err
	}
	//fill response
	response := &models.TokenResponse{
		Access:  accessToken,
		Refresh: refreshToken,
	}
	return response, nil
}
