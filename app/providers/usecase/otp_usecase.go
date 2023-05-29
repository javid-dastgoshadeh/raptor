package usecase

import (
	"context"
	"fmt"
	"raptor/app/providers"
	"raptor/cache"
	env "raptor/config"
	"raptor/models"
	"raptor/pkg/helpers"
	"strings"
)

// NewOtp ...
type otpUsecase struct {
	KratosRepo providers.KratosRepository
	OtpRepo    providers.OtpRepository
}

// NewOtp ...
func NewOtp(kratosRepo providers.KratosRepository, otpRepo providers.OtpRepository) *otpUsecase {
	return &otpUsecase{
		KratosRepo: kratosRepo,
		OtpRepo:    otpRepo,
	}
}

// SendOtpCode ...
func (u *otpUsecase) SendOtpCode(ctx context.Context, emailMobile string) error {

	//generate random sms code
	code := helpers.GenerateRandomInt()
	//check identity format it must be email or mobile
	identityFormat, err := helpers.CheckFormat(emailMobile)
	if err != nil {
		return err
	}

	//check otp code sender type(email or sms)
	senderType := models.MessageSender(identityFormat)

	//cache code for identity
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheOtpCode(emailMobile), code)
	//get timestamp for otp expire time
	expireTime := helpers.GenerateTimeStampForFutureBaseOnMinute(env.GetInt("courier.sms.expire_time"))
	//cache time for code sent
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheTimeForOtpCode(emailMobile), fmt.Sprintf("%v", expireTime))

	//check user existence
	ok, err := u.KratosRepo.CheckIdentityExistence(ctx, emailMobile, string(senderType))

	if err != nil {
		return err
	}

	//set flowID in this variable
	var flow string

	//if user is existed we create a login flow and if is not existed we create register flow
	if ok {
		//create a login flow on kratos
		flow, err = u.KratosRepo.CreateKratosLoginFLow(ctx)
		if err != nil {
			return err
		}
		//cache request state(register or login)
		cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheState(emailMobile), "login")
	}

	if !ok {
		//create a register flow on kratos
		flow, err = u.KratosRepo.CreateKratosRegisterFLow(ctx)
		if err != nil {
			return err
		}
		//cache request state(register or login)
		cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheState(emailMobile), "register")
	}
	//cache flowId for identity
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheFlowID(emailMobile), flow)

	err = u.OtpRepo.SendCode(ctx, emailMobile, code, &senderType)
	if err != nil {
		return err
	}
	return nil
}

// VerifyOtpCode ...
func (u *otpUsecase) VerifyOtpCode(ctx context.Context, emailMobile string, code string, device string) (*models.TokenResponse, error) {

	var (
		err error
		res interface{}
	)
	//read code from cache by identity
	err = u.OtpRepo.VerifyCode(ctx, emailMobile, code)
	if err != nil {
		return nil, err
	}

	//read flow from cache by identity
	cachedFlow := cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheFlowID(emailMobile))
	//read state from cache by identity
	cachedState := cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheState(emailMobile))

	identityFormat, err := helpers.CheckFormat(emailMobile)
	if err != nil {
		return nil, err
	}

	//submit login flow for login state
	if cachedState == "login" {
		submitLoginReq := &models.SubmitKratosLoginRequest{
			Method:     "password",
			Identifier: emailMobile,
			Password:   env.GetString("identity.password"),
		}
		res, err = u.KratosRepo.SubmitKratosLoginFlow(ctx, cachedFlow, submitLoginReq)
		if err != nil {
			return nil, err
		}
	}

	//submit register flow for login register
	if cachedState == "register" {

		// if device type is app we ask user to set display name then response jwt to it
		var deviceType models.Device
		if device == "" || device == string(models.Web) {
			deviceType = models.Web
		}
		if device == string(models.App) {
			deviceType = models.App
		}
		//traits is a structure for user attributes
		var traits *models.Traits
		sms := models.Sms
		email := models.Email

		var (
			randomUsername string
			isExist        = true
		)
		//check is username exist before and if existed generate another username
		for isExist {
			randomUsername = strings.ToLower(helpers.GenerateRandomString(8))
			isExist, err = u.KratosRepo.CheckIdentityExistence(ctx, randomUsername, string(models.Username))
			if err != nil {
				return nil, err
			}
		}

		//check identity format
		//identity can be in two format sms and email
		if identityFormat == string(email) {
			traits = &models.Traits{
				Email:         emailMobile,
				DisplayName:   env.GetString("identity.default_display_name"),
				EmailVerified: "true",
				Username:      randomUsername,
			}
		}

		if identityFormat == string(sms) {
			traits = &models.Traits{
				PhoneNumber:         emailMobile,
				DisplayName:         env.GetString("identity.default_display_name"),
				PhoneNumberVerified: "true",
				Username:            randomUsername,
			}
		}

		submitLoginReq := &models.SubmitKratosRegisterRequest{
			Method:   "password",
			Traits:   traits,
			Password: env.GetString("identity.password"),
		}
		res, err = u.KratosRepo.SubmitKratosRegisterFLow(ctx, cachedFlow, submitLoginReq)

		if err != nil {
			return nil, err
		}

		if deviceType == models.App {
			return nil, nil
		}
	}

	//generate access_token
	//parse response to get desire claims for jwt token
	claims, err := helpers.GenerateAccessTokenClaims(res)
	if err != nil {
		return nil, err
	}

	//check account that must be active
	_, err = u.KratosRepo.ActiveIdentity(ctx, claims, identityFormat)
	if err != nil {
		return nil, err
	}

	//check updated tokens
	claims, err = helpers.GenerateAccessTokenClaims(res)
	if err != nil {
		return nil, err
	}

	//load private key for access_token
	//privateKey := helpers.LoadPrivateKey(env.GetString("security.private_key_path"))
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
	//refreshTokenPrivateKey := helpers.LoadPrivateKey("keys/refresh_token_private_key.pem")
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

// AfterRegisterStep ...
func (u *otpUsecase) AfterRegisterStep(ctx context.Context, emailMobile string, code string, displayName string) (*models.TokenResponse, error) {

	var (
		err error
		res interface{}
	)
	//read code from cache by identity
	err = u.OtpRepo.VerifyCode(ctx, emailMobile, code)
	if err != nil {
		return nil, err
	}

	//create a login flow on kratos
	flow, err := u.KratosRepo.CreateKratosLoginFLow(ctx)
	if err != nil {
		return nil, err
	}

	submitLoginReq := &models.SubmitKratosLoginRequest{
		Method:     "password",
		Identifier: emailMobile,
		Password:   env.GetString("identity.password"),
	}
	res, err = u.KratosRepo.SubmitKratosLoginFlow(ctx, flow, submitLoginReq)
	if err != nil {
		return nil, err
	}

	injection := make(map[string]interface{})
	injection["display_name"] = displayName
	updatedRes, err := helpers.InjectDataToIdentity(res, injection)

	if err != nil {
		return nil, err
	}

	sessionInfo, ok := res.(map[string]interface{})
	if !ok {
		return nil, models.ErrCreatingRegisterFlow
	}
	dataInfo := sessionInfo["session"].(map[string]interface{})
	identityInfo := dataInfo["identity"].(map[string]interface{})

	//check if state is register add update use and add public metadata 1 for show user activate
	_, err = u.KratosRepo.UpdateIdentity(ctx, fmt.Sprintf("%v", identityInfo["id"]), updatedRes)
	if err != nil {
		return nil, err
	}

	//generate access_token
	//parse response to get desire claims for jwt token
	claims, err := helpers.GenerateAccessTokenClaims(sessionInfo)
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
