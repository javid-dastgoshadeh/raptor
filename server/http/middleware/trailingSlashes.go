package middleware

import echoMiddleware "github.com/labstack/echo/v4/middleware"

// RegisterTrailingSlashes ...
func RegisterTrailingSlashes() {
	Server.Pre(echoMiddleware.RemoveTrailingSlash())
}
