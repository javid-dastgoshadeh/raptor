package middleware

import (
	"github.com/labstack/echo/v4"

	"raptor/config"
)

// Server instance
var Server *echo.Echo

// Init ...
func Init(server *echo.Echo) {
	Server = server
	if !env.GetBool("debug") {
		RegisterCSRF()
		RegisterLogger()
		RegisterRateLimit()
	}
	RegisterRecover()
	RegisterTrailingSlashes()
	RegisterCORS()

}
