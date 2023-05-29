package http

import (
	"github.com/labstack/echo/v4"
	"raptor/app/providers/repository"
	"raptor/app/providers/repository/methods"
	"raptor/app/providers/transport/http"
	"raptor/app/providers/usecase"
)

// Register ...
func Register(e *echo.Echo) {

	var email = methods.Email{}
	var sms = methods.Sms{}
	otpRepo := repository.OtpRepo(email, sms)
	socialRepo := repository.SocialRepo()
	KratosRepo := repository.KratosRepo()
	ou := usecase.NewOtp(KratosRepo, otpRepo)
	su := usecase.NewSocial(KratosRepo, socialRepo)
	tu := usecase.NewToken(KratosRepo)
	http.RegisterHandlers(e, ou, su, tu)
}
