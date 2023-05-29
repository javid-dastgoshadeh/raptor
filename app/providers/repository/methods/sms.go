package methods

import (
	"context"
	"raptor/pkg/helpers"
)

// Sms ...
type Sms struct{}

func (sms *Sms) SendCode(ctx context.Context, UserIdentification string, code string) error {
	go func() {
		helpers.KavenegarSmsSender(code, UserIdentification)
	}()
	return nil
}
