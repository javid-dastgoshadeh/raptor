package http

import (
	"fmt"
	"raptor/cache"
	"raptor/models"
	"raptor/pkg/helpers"
)

func CheckCodeSentInterval(identity string) error {
	if cache.RedisInstance.GetValue(helpers.GenerateIdentityKeyToCacheTimeForOtpCode(identity)) > fmt.Sprintf("%v", helpers.GenerateCurrentTimeStamp()) {
		return models.ErrCodeSentBefore
	}
	return nil
}

func CheckIdentityValidFormat(identity string) error {
	if !helpers.CheckEmailFormat(identity) && !helpers.CheckMobileFormat(identity) {
		return models.ErrIdentityFormat
	}
	return nil
}
