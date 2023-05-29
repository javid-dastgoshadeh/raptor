package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// RegisterCSRF ...
func RegisterCSRF() {
	Server.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		TokenLength:  32,
		TokenLookup:  "header:" + echo.HeaderXCSRFToken,
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}))
}
