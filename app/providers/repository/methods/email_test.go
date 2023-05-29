package methods

import (
	"context"
	env "raptor/config"
	"testing"
)

func TestEmail_SendCode(t *testing.T) {
	email := Email{}
	env.Init("./config_sample.json")
	ctx := context.Background()
	email.SendCode(ctx, "javiddastgoshadeh@gmail.com", "12345")
}
