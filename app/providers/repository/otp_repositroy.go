package repository

import (
	"context"
	"fmt"
	"raptor/app/providers"
	"raptor/logger"
	"raptor/pkg/helpers"

	"raptor/app/providers/repository/methods"
	"raptor/cache"
	"raptor/models"
)

type OtpRepository struct {
	email methods.Email
	sms   methods.Sms
}

// OtpRepo ...
func OtpRepo(email methods.Email, sms methods.Sms) providers.OtpRepository {
	return &OtpRepository{
		email: email,
		sms:   sms,
	}
}

// SendCode ...
func (repo *OtpRepository) SendCode(ctx context.Context, UserIdentification string, code string, senderType *models.MessageSender) error {

	if *senderType == models.Sms {
		return repo.sms.SendCode(ctx, UserIdentification, code)
	}

	if *senderType == models.Email {
		return repo.email.SendCode(ctx, UserIdentification, code)
	}
	return nil
}

// VerifyCode ...
func (repo *OtpRepository) VerifyCode(_ context.Context, identity string, code string) error {
	if cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheTimeForOtpCode(identity)) < fmt.Sprintf("%v", helpers.GenerateCurrentTimeStamp()) {
		logger.Error(models.ErrExpireOtpCodeTime.Error())
		return models.ErrExpireOtpCodeTime
	}
	if cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheOtpCode(identity)) == code {
		logger.Error(models.ErrOtpCodeNotMatched)
		return nil
	}
	return models.ErrOtpCodeNotMatched
}
