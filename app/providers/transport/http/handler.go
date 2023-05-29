package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"raptor/app/providers"
	"raptor/config"
	"raptor/logger"
	"raptor/models"
	"raptor/pkg/helpers"
	"raptor/pkg/templates"
)

// Handler
type handler struct {
	UC providers.OtpUsecase
	SU providers.SocialUsecase
	TU providers.TokenUsecase
}

// SendCode ...
func (h *handler) SendCode(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	request := &models.CodeRequest{}
	context := ctx.Request().Context()
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(models.ErrIdentityFormat.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}

	if err = ctx.Validate(request); err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}
	err = CheckIdentityValidFormat(request.EmailMobile)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}
	err = CheckCodeSentInterval(request.EmailMobile)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	err = h.UC.SendOtpCode(context, request.EmailMobile)

	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok("code send successfully", nil)
	return ctx.JSON(http.StatusOK, output)
}

// VerifyCode ...
func (h *handler) VerifyCode(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	request := &models.VerifyCodeRequest{}
	context := ctx.Request().Context()
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(models.ErrIdentityFormat.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}

	if err = ctx.Validate(request); err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}
	err = CheckIdentityValidFormat(request.EmailMobile)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	response, err := h.UC.VerifyOtpCode(context, request.EmailMobile, request.Code, request.Device)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	//TODO in next version
	//this case for response for register step
	if err == nil && response == nil {
		resMsg := make(map[string]interface{})
		resMsg["register"] = 1
		template = templates.MobileAppRegisterResponse(resMsg)
		return ctx.JSON(http.StatusOK, template)
	}
	output := templates.Ok(response, nil)
	return ctx.JSON(http.StatusOK, output)
}

// AfterRegisterStep ...
func (h *handler) AfterRegisterStep(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	request := &models.AfterRegisterRequest{}

	context := ctx.Request().Context()

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(models.ErrIdentityFormat.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}

	if err = ctx.Validate(request); err != nil {
		logger.Error(err.Error())
		template = templates.BadRequest(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	response, err := h.UC.AfterRegisterStep(context, request.EmailMobile, request.Code, request.Name)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(response, nil)
	return ctx.JSON(http.StatusOK, output)
}

// SocialVerify ...
func (h *handler) SocialVerify(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	msg := errors.New("internal server error")
	var (
		payload map[string]interface{}
		err     error
	)
	context := ctx.Request().Context()

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	provider := ctx.QueryParam("provider")

	err = json.Unmarshal(body, &payload)

	response, err := h.SU.VerifySocialCode(context, provider, payload)
	if err != nil {
		if env.GetBool("debug") {
			msg = err
		}
		template = templates.MobileAppInternalServerError(msg.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(response, nil)
	return ctx.JSON(http.StatusOK, output)
}

// ExchangeSocialCode ...
func (h *handler) ExchangeSocialCode(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	//msg := errors.New("internal server error")
	var (
		payload map[string]interface{}
		err     error
	)
	context := ctx.Request().Context()

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	provider := ctx.QueryParam("provider")

	err = json.Unmarshal(body, &payload)

	response, err := h.SU.ExchangeSocialCode(context, provider, payload)
	if err != nil {

		template = templates.MobileAppInternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(response, nil)
	return ctx.JSON(http.StatusOK, output)
}

// RefreshToken ...
func (h *handler) RefreshToken(ctx echo.Context) error {

	c := ctx.Request().Context()

	publicKey := models.RefreshTokenPublicKey

	tokenString, err := helpers.GetJwtStringFromRequest(ctx.Request())
	if err != nil {
		templates.Unauthorized(err)
		return ctx.JSON(http.StatusUnauthorized, templates.Unauthorized(err))
	}
	// Retrieve JWT from context
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		logger.Error("err")
		return ctx.JSON(http.StatusUnauthorized, templates.Forbidden(err))
	}
	// Access claims
	claims := token.Claims.(jwt.MapClaims)

	res, err := h.TU.RefreshToken(c, claims)
	if err != nil {
		logger.Error(err.Error())
		template := templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(res, nil)
	return ctx.JSON(http.StatusOK, output)
}

// RegisterHandlers ...
func RegisterHandlers(e *echo.Echo, ou providers.OtpUsecase, su providers.SocialUsecase, tu providers.TokenUsecase) {
	handler := &handler{
		UC: ou,
		SU: su,
		TU: tu,
	}
	router := e.Group("/api/v1")
	{
		router.POST("/send-code", handler.SendCode)
		router.POST("/verify-code", handler.VerifyCode)
		router.POST("/update-username", handler.AfterRegisterStep)
		router.POST("/verify-social", handler.SocialVerify)
		router.POST("/exchange-social-code", handler.ExchangeSocialCode)
		router.GET("/refresh-token", handler.RefreshToken)
	}
}
