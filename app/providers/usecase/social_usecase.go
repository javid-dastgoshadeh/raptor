package usecase

import (
	"context"
	"fmt"
	"raptor/app/providers"
	"raptor/cache"
	env "raptor/config"
	"raptor/logger"
	"raptor/models"
	"raptor/pkg/helpers"
)

// socialUsecase ...
type socialUsecase struct {
	KratosRepo providers.KratosRepository
	SocialRepo providers.SocialRepository
}

// NewSocial ...
func NewSocial(kratosRepo providers.KratosRepository, socialRepo providers.SocialRepository) *socialUsecase {
	return &socialUsecase{
		KratosRepo: kratosRepo,
		SocialRepo: socialRepo,
	}
}

// VerifySocialCode ...
func (u *socialUsecase) VerifySocialCode(ctx context.Context, provider string, data map[string]interface{}) (*models.TokenResponse, error) {

	var (
		response *models.TokenResponse
		err      error
		identity *models.Traits
		flow     string
		res      interface{}
	)
	//check providers
	switch provider {

	case "google":
		identity, err = u.SocialRepo.VerifyGoogle(ctx, data)
		fmt.Println("google", data)
	case "apple":
		identity, err = u.SocialRepo.VerifyApple(ctx, data)
		fmt.Println("apple", data)
	}

	if err != nil {
		return nil, err
	}

	//check user existence
	ok, err := u.KratosRepo.CheckIdentityExistence(ctx, fmt.Sprintf("%v", identity.Email), string(models.Email))

	if ok {
		//create a login flow on kratos
		flow, err = u.KratosRepo.CreateKratosLoginFLow(ctx)
		if err != nil {
			return nil, err
		}
		//cache request state(register or login)
		cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheState(fmt.Sprintf("%v", identity.Email)), "login")
	}

	if !ok {
		//create a register flow on kratos
		flow, err = u.KratosRepo.CreateKratosRegisterFLow(ctx)
		if err != nil {
			return nil, err
		}
		//cache request state(register or login)
		cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheState(fmt.Sprintf("%v", identity.Email)), "register")
	}

	//cache flow id for identity
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheFlowID(fmt.Sprintf("%v", identity.Email)), flow)

	//read code from cache by identity
	cachedFlow := cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheFlowID(fmt.Sprintf("%v", identity.Email)))
	//read state from cache by identity
	cachedState := cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheState(fmt.Sprintf("%v", identity.Email)))

	//submit login flow for login state
	if cachedState == "login" {
		submitLoginReq := &models.SubmitKratosLoginRequest{
			Method:     "password",
			Identifier: fmt.Sprintf("%v", identity.Email),
			Password:   env.GetString("identity.password"),
		}
		res, err = u.KratosRepo.SubmitKratosLoginFlow(ctx, cachedFlow, submitLoginReq)

	}

	//submit register flow for login register
	if cachedState == "register" {

		//traits is a structure for user attributes
		var traits *models.Traits

		traits = &models.Traits{
			PhoneNumber: identity.Email,
		}
		submitLoginReq := &models.SubmitKratosRegisterRequest{
			Method:   "password",
			Traits:   traits,
			Password: env.GetString("identity.password"),
		}
		res, err = u.KratosRepo.SubmitKratosRegisterFLow(ctx, cachedFlow, submitLoginReq)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	//parse response to get desire claims for jwt token
	claims, err := helpers.GenerateAccessTokenClaims(res)
	if err != nil {
		return nil, err
	}
	//load private key for access_token
	privateKey := helpers.LoadPrivateKey(env.GetString("security.private_key_path"))
	accessToken, err := helpers.GenerateJwtToken(privateKey, claims)
	if err != nil {
		return nil, err
	}

	//parse response to get desire claims for jwt token
	refreshTokenClaims, err := helpers.GenerateRefreshTokenClaims(res)
	if err != nil {
		return nil, err
	}
	//load private key for refresh_token
	refreshTokenPrivateKey := helpers.LoadPrivateKey(env.GetString("security.refresh_token_private_key_path"))

	refreshToken, err := helpers.GenerateJwtToken(refreshTokenPrivateKey, refreshTokenClaims)
	if err != nil {
		return nil, err
	}
	//fill response
	response = &models.TokenResponse{Access: accessToken, Refresh: refreshToken}
	return response, nil

}

// ExchangeSocialCode ...
func (u *socialUsecase) ExchangeSocialCode(ctx context.Context, provider string, data map[string]interface{}) (*models.TokenResponse, error) {

	var (
		response *models.TokenResponse
		err      error
		identity *models.Traits
		flow     string
		res      interface{}
	)
	//check providers
	switch provider {

	case "google":
		identity, err = u.SocialRepo.ExchangeGoogleOauth2Code(ctx, data)
		fmt.Println("google", data)
		//case "apple":
		//	identity, err = u.SocialRepo.ExchangeAppleOauth2Code(ctx, data)
		//	fmt.Println("apple", data)
	}

	if err != nil {
		return nil, err
	}

	//check user existence
	ok, err := u.KratosRepo.CheckIdentityExistence(ctx, fmt.Sprintf("%v", identity.Email), string(models.Email))

	if ok {
		//create a login flow on kratos
		flow, err = u.KratosRepo.CreateKratosLoginFLow(ctx)
		if err != nil {
			return nil, err
		}
		//cache request state(register or login)
		cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheState(fmt.Sprintf("%v", identity.Email)), "login")
	}

	if !ok {
		//create a register flow on kratos
		flow, err = u.KratosRepo.CreateKratosRegisterFLow(ctx)
		if err != nil {
			return nil, err
		}
		//cache request state(register or login)
		cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheState(fmt.Sprintf("%v", identity.Email)), "register")
	}

	//cache flow id for identity
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheFlowID(fmt.Sprintf("%v", identity.Email)), flow)

	//read flow id from cache by identity
	cachedFlow := cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheFlowID(fmt.Sprintf("%v", identity.Email)))
	//read state from cache by identity
	cachedState := cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheState(fmt.Sprintf("%v", identity.Email)))

	//submit login flow for login state
	if cachedState == "login" {
		submitLoginReq := &models.SubmitKratosLoginRequest{
			Method:     "password",
			Identifier: fmt.Sprintf("%v", identity.Email),
			Password:   env.GetString("identity.password"),
		}
		res, err = u.KratosRepo.SubmitKratosLoginFlow(ctx, cachedFlow, submitLoginReq)
	}

	//submit register flow for login register
	if cachedState == "register" {

		submitLoginReq := &models.SubmitKratosRegisterRequest{
			Method:   "password",
			Traits:   identity,
			Password: env.GetString("identity.password"),
		}
		res, err = u.KratosRepo.SubmitKratosRegisterFLow(ctx, cachedFlow, submitLoginReq)
	}

	//parse response to get desire claims for jwt token
	claims, err := helpers.GenerateAccessTokenClaims(res)
	if err != nil {
		return nil, err
	}
	//load private key for access_token
	privateKey := helpers.LoadPrivateKey(env.GetString("security.private_key_path"))
	accessToken, err := helpers.GenerateJwtToken(privateKey, claims)
	if err != nil {
		return nil, err
	}

	//upload avatar
	if cachedState == "register" && identity.Avatar != nil {
		err := helpers.UploadImage(accessToken, fmt.Sprintf("%v", identity.Avatar))
		if err != nil {
			logger.Error(err)
		}
	}

	//parse response to get desire claims for jwt token
	refreshTokenClaims, err := helpers.GenerateRefreshTokenClaims(res)
	if err != nil {
		return nil, err
	}
	//load private key for refresh_token
	refreshTokenPrivateKey := helpers.LoadPrivateKey(env.GetString("security.refresh_token_private_key_path"))

	refreshToken, err := helpers.GenerateJwtToken(refreshTokenPrivateKey, refreshTokenClaims)
	if err != nil {
		return nil, err
	}
	//fill response
	response = &models.TokenResponse{Access: accessToken, Refresh: refreshToken}
	return response, nil

}
