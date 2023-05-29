package methods

import (
	"context"
	"testing"

	env "raptor/config"
)

func TestSms_SendCode(t *testing.T) {
	email := Sms{}
	env.Init("./config.json")
	ctx := context.Background()
	email.SendCode(ctx, "09125317278", "2555")
}
