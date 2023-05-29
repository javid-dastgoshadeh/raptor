package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"raptor/app/account"
	"raptor/logger"
	"raptor/models"
	"raptor/pkg/helpers"
	"raptor/pkg/templates"
	"raptor/server/http/middleware"
)

// Handler
type handler struct {
	UC account.Usecase
}

// UpdateProfile ...
func (h *handler) UpdateProfile(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	context := ctx.Request().Context()

	var (
		err       error
		updateReq *models.UpdateRequest
	)

	claims, err := helpers.ExtractClaimsFromToken(ctx)

	if err != nil {
		logger.Error(err.Error())
		template = templates.Unauthorized(err.Error())
		return ctx.JSON(http.StatusUnauthorized, template)
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	err = json.Unmarshal(body, &updateReq)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	res, err := h.UC.UpdateProfile(context, claims, updateReq)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(res, nil)
	return ctx.JSON(http.StatusOK, output)
}
func (h *handler) UpdateIdentifier(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	context := ctx.Request().Context()
	var (
		payload map[string]interface{}
		err     error
	)

	tokenString, err := helpers.GetJwtStringFromRequest(ctx.Request())

	//claims, err := helpers.ExtractClaimsFromToken(ctx)

	if err != nil {
		logger.Error(err.Error())
		template = templates.Unauthorized(err.Error())
		return ctx.JSON(http.StatusUnauthorized, template)
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	msg, err := h.UC.UpdateIdentifier(context, tokenString, payload)
	//TODO Error handeling
	if err == models.ErrUpdateMobile {
		output := templates.Ok(msg, nil)
		return ctx.JSON(http.StatusRequestedRangeNotSatisfiable, output)
	}

	if err != nil {
		logger.Error(err.Error())
		template = templates.MobileAppInternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(msg, nil)
	return ctx.JSON(http.StatusOK, output)
}
func (h *handler) VerifyUpdateIdentifier(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	context := ctx.Request().Context()

	var (
		payload map[string]interface{}
		err     error
	)

	claims, err := helpers.ExtractClaimsFromToken(ctx)
	if err != nil {
		logger.Error(err.Error())
		template = templates.Unauthorized(err.Error())
		return ctx.JSON(http.StatusUnauthorized, template)
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusBadRequest, template)
	}

	res, err := h.UC.VerifyUpdateIdentifier(context, claims, payload)

	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(res, nil)
	return ctx.JSON(http.StatusOK, output)
}
func (h *handler) Profile(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	context := ctx.Request().Context()
	var (
		res interface{}
		err error
	)

	claims, err := helpers.ExtractClaimsFromToken(ctx)
	if err != nil {
		logger.Error(err.Error())
		template = templates.Unauthorized(err.Error())
		return ctx.JSON(http.StatusUnauthorized, template)
	}

	res, err = h.UC.Profile(context, claims)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}
	output := templates.Ok(res, nil)
	return ctx.JSON(http.StatusOK, output)
}
func (h *handler) Logout(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	var err error
	context := ctx.Request().Context()

	claims, err := helpers.ExtractClaimsFromToken(ctx)
	if err != nil {
		logger.Error(err.Error())
		template = templates.Unauthorized(err.Error())
		return ctx.JSON(http.StatusUnauthorized, template)
	}

	err = h.UC.Logout(context, claims)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}

	output := templates.Ok("logout successfully", nil)
	return ctx.JSON(http.StatusOK, output)
}
func (h *handler) InactiveAccount(ctx echo.Context) error {
	template := &templates.ResponseTemplate{}
	context := ctx.Request().Context()

	var err error

	claims, err := helpers.ExtractClaimsFromToken(ctx)
	if err != nil {
		logger.Error(err.Error())
		template = templates.Unauthorized(err.Error())
		return ctx.JSON(http.StatusUnauthorized, template)
	}

	err = h.UC.InactivateAccount(context, claims)
	if err != nil {
		logger.Error(err.Error())
		template = templates.InternalServerError(err.Error())
		return ctx.JSON(http.StatusInternalServerError, template)
	}

	output := templates.Ok("inactive successfully", nil)
	return ctx.JSON(http.StatusOK, output)
}

func RegisterHandlers(e *echo.Echo, uc account.Usecase) {
	//add authentication middleware
	authenticationMid := middleware.RegisterAuthentication()

	handler := &handler{UC: uc}
	router := e.Group("/api/v1")

	{
		//add authentication middleware to this group
		router.Use(authenticationMid)
		router.POST("/update-profile", handler.UpdateProfile)
		//router.POST("/verification", handler.Verification)
		//router.POST("/verify-verification", handler.VerifyVerification)
		router.POST("/update-identifier", handler.UpdateIdentifier)
		router.POST("/verify-update-identifier", handler.VerifyUpdateIdentifier)
		router.GET("/profile", handler.Profile)
		router.GET("/logout", handler.Logout)
		router.GET("/inactive-account", handler.InactiveAccount)
	}
}
