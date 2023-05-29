package http

import (
	"github.com/labstack/echo/v4"
	"raptor/app/providers/repository/methods"

	"raptor/app/account/repository"
	"raptor/app/account/transport/http"
	"raptor/app/account/usecase"
	otpRepository "raptor/app/providers/repository"
)

// Register ...
func Register(e *echo.Echo) {
	var email = methods.Email{}
	var sms = methods.Sms{}
	repo := repository.KratosRepo()
	identificationRepo := repository.IdentificationRepo()
	otpRepo := otpRepository.OtpRepo(email, sms)
	u := usecase.Usecase(repo, otpRepo, identificationRepo)
	http.RegisterHandlers(e, u)
}
