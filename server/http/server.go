package httpServer

import (
	"net/http"
	"os"
	"raptor/models"
	"raptor/pkg/helpers"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"

	env "raptor/config"
	"raptor/logger"
	"raptor/pkg/templates"
	"raptor/server/http/gateway"
	"raptor/server/http/middleware"
)

// instance is main echo server instance
var instance *echo.Echo

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Debug   bool   `json:"debug"`
}

// HTTPErrorHandler ...
func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
	}
	logger.Error(err)
	//c.Logger().Error(err)
	template := templates.GetWithCode(code, err)
	c.JSON(code, template)
}

// Serve create new echo server and run it
func Serve(cfg *ServerConfig) {
	instance = echo.New()
	instance.HideBanner = true
	if env.GetBool("debug") {
		instance.Debug = cfg.Debug
	}
	//add internal error template to response
	instance.HTTPErrorHandler = HTTPErrorHandler
	instance.Validator = &Validator{validator: validator.New()}

	// Register middlewares
	middleware.Init(instance)

	//load privateKey
	models.PrivateKey = helpers.LoadPrivateKey(env.GetString("security.private_key_path"))
	//load refreshTokenPrivateKey
	models.RefreshTokenPrivateKey = helpers.LoadPrivateKey(env.GetString("security.refresh_token_private_key_path"))
	//load RefreshTokenPublicKey
	models.RefreshTokenPublicKey = helpers.LoadPublicKey(env.GetString("security.refresh_token_public_key_path"))

	// Register services
	gateway.Init(instance)
	//instance.Debug = cfg.Debug
	logger.Info("Server running at " + cfg.Address + ":" + strconv.Itoa(cfg.Port))
	//instance.Logger.Info("Server running at " + cfg.Address + ":" + strconv.Itoa(cfg.Port))

	//logger.Info("Server running at " + cfg.Address + ":" + strconv.Itoa(cfg.Port))

	err := instance.Start(cfg.Address + ":" + strconv.Itoa(cfg.Port))

	if err != nil {
		instance.Logger.Fatal(err)

		os.Exit(1)
	}
}

// GetInstance return main server instance
func GetInstance() *echo.Echo {
	return instance
}
