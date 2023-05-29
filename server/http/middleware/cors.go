package middleware

import (
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// RegisterCORS ...
func RegisterCORS() {
	Server.Use(echoMiddleware.CORS())
}
