package middleware

import "github.com/labstack/echo/v4/middleware"

func RegisterRecover() {
	Server.Use(
		middleware.RecoverWithConfig(
			middleware.DefaultRecoverConfig,
		),
	)

}
