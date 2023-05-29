package methods

import (
	"context"
	env "raptor/config"
	"raptor/logger"
	"raptor/pkg/helpers"
)

// Email ...
type Email struct{}

func (email *Email) SendCode(ctx context.Context, UserIdentification string, code string) error {
	go func() {
		preCode := env.GetString("templates.email_verification_code.content.before_code")
		afterCode := env.GetString("templates.email_verification_code.content.after_code")
		msg := preCode + code + afterCode
		helpers.SendMail(
			UserIdentification,
			env.GetString("templates.email_verification_code.subject"),
			msg,
			env.GetString("templates.email_verification_code.template_name"),
		)
		logger.Info("this code: " + code + " is sent to " + UserIdentification)
	}()
	return nil
}
