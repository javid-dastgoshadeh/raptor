package gateway

import (
	"github.com/labstack/echo/v4"
	accountProvider "raptor/app/account/registry/http"
	otpProvider "raptor/app/providers/registry/http"
)

// Init initialize http services
func Init(server *echo.Echo) {
	otpProvider.Register(server)
	accountProvider.Register(server)
}
