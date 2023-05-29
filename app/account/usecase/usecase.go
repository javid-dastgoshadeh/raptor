package usecase

import (
	"context"
	"fmt"
	"raptor/app/account"
	"raptor/app/providers"
	"raptor/cache"
	env "raptor/config"
	"raptor/models"
	"raptor/pkg/helpers"
)

// Usecase ...
type usecase struct {
	Repo    account.KratosRepository
	OtpRepo providers.OtpRepository
	IdeRepo account.IdentificationRepository
}

// Usecase ...
func Usecase(repo account.KratosRepository, otpRepo providers.OtpRepository, ideRepo account.IdentificationRepository) *usecase {
	return &usecase{
		Repo:    repo,
		OtpRepo: otpRepo,
		IdeRepo: ideRepo,
	}
}

// Profile ...
func (u *usecase) Profile(ctx context.Context, jwtClaims map[string]interface{}) (interface{}, error) {
	res, err := u.Repo.CheckSession(ctx, jwtClaims)
	if err != nil {
		return nil, err
	}
	result, err := helpers.ExportIdentityFromKratosResponse(res)
	if err != nil {
		return nil, err
	}
	return result, err
}

// UpdateProfile ...
func (u *usecase) UpdateProfile(ctx context.Context, jwtClaims map[string]interface{}, request *models.UpdateRequest) (interface{}, error) {
	//create setting flow to update profile
	flow, err := u.Repo.CreateSettingFlow(ctx, jwtClaims)
	//traits := request.Traits

	//fill identity from update profile request
	identity := &models.SubmitKratosSettingRequest{
		Method: "profile",
		Traits: request.Traits,
	}

	//get identity by id from admin side
	getIdentity, err := u.Repo.GetIdentity(ctx, jwtClaims)
	if err != nil {
		return nil, err
	}
	//export traits from getIdentity response
	exportTraits, err := helpers.GetTraitsFromKratosResponse(getIdentity)

	//check if identifier is not changed in request
	if identity.Traits.Email != exportTraits.Email || identity.Traits.PhoneNumber != exportTraits.PhoneNumber {
		return nil, models.ErrChangeIdentifier
	}
	identity.Traits.EmailVerified = exportTraits.EmailVerified
	identity.Traits.PhoneNumberVerified = exportTraits.PhoneNumberVerified
	res, err := u.Repo.SubmitSettingFlow(ctx, identity, flow, jwtClaims)
	if err != nil {
		return nil, err
	}

	result, err := helpers.ExportIdentityFromKratosResponse(res)
	if err != nil {
		return nil, err
	}

	return result, err
}

// UpdateIdentifier ...
func (u *usecase) UpdateIdentifier(ctx context.Context, jwt string, identity map[string]interface{}) (interface{}, error) {
	//generate random sms code
	code := helpers.GenerateRandomInt()
	//check identity format it must be email or mobile
	identityFormat, err := helpers.CheckFormat(fmt.Sprintf("%v", identity["identity"]))
	if err != nil {
		return "", err
	}

	//check otp code sender type(email or sms)
	senderType := models.MessageSender(identityFormat)
	//check if identity is mobile number and verify before,can't be changed
	if senderType == models.Sms {
		ok, err := u.IdeRepo.CheckIdentification(ctx, jwt)
		if err != nil {
			return "", err
		}
		if ok {
			return "you can't update mobile number after verification", models.ErrUpdateMobile
		}
	}
	//cache code for identity
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheOtpCode(fmt.Sprintf("%v", identity["identity"])), code)
	//get timestamp for otp expire time
	expireTime := helpers.GenerateTimeStampForFutureBaseOnMinute(env.GetInt("courier.sms.expire_time"))
	//cache time for code sent
	cache.RedisInstance.SetValue(helpers.GenerateIdentityKeyToCacheTimeForOtpCode(fmt.Sprintf("%v", identity["identity"])), fmt.Sprintf("%v", expireTime))

	err = u.OtpRepo.SendCode(ctx, fmt.Sprintf("%v", identity["identity"]), code, &senderType)

	if err != nil {
		return "", err
	}

	return "code send successfully to new identity", nil
}

// VerifyUpdateIdentifier ...
func (u *usecase) VerifyUpdateIdentifier(ctx context.Context, jwtClaims map[string]interface{}, identity map[string]interface{}) (interface{}, error) {

	var (
		err error
		res interface{}
	)
	//read code from cache by identity
	err = u.OtpRepo.VerifyCode(ctx, fmt.Sprintf("%v", identity["identity"]), fmt.Sprintf("%v", identity["code"]))
	if err != nil {
		return nil, err
	}

	identityFormat, err := helpers.CheckFormat(fmt.Sprintf("%v", identity["identity"]))
	if err != nil {
		return nil, err
	}

	//get identity by id from admin side
	getIdentity, err := u.Repo.GetIdentity(ctx, jwtClaims)
	if err != nil {
		return nil, err
	}

	injection := make(map[string]interface{})

	if identityFormat == "email" {
		injection["email"] = identity["identity"]
		injection["email_verified"] = "true"
	}
	if identityFormat == "sms" {
		injection["phone_number"] = identity["identity"]
		injection["phone_number_verified"] = "true"
	}

	getIdentity, err = helpers.InjectDataToIdentity(getIdentity, injection)
	if err != nil {
		return nil, err
	}
	res, err = u.Repo.UpdateIdentity(ctx, fmt.Sprintf("%v", jwtClaims["sub"]), getIdentity)

	if err != nil {
		return nil, err
	}

	result, err := helpers.ExportIdentityFromKratosResponse(res)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// InactivateAccount ...
func (u *usecase) InactivateAccount(ctx context.Context, jwtClaims map[string]interface{}) error {
	err := u.Repo.InactiveIdentity(ctx, jwtClaims)
	return err
}

// Logout ...
func (u *usecase) Logout(ctx context.Context, jwtClaims map[string]interface{}) error {
	err := u.Repo.DisableSession(ctx, jwtClaims)
	return err
}
