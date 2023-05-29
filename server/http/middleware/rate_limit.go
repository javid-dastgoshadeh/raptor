package middleware

import (
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	env "raptor/config"
)

// RegisterRateLimit ...
func RegisterRateLimit() {
	Server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(env.GetInt("services.rate_limit")))))
}
