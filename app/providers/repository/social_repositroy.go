package repository

import (
	"context"
	"fmt"
	"github.com/Timothylock/go-signin-with-apple/apple"
	"raptor/app/providers"
	env "raptor/config"
	"raptor/logger"
	"raptor/models"
	"raptor/pkg/helpers"
	"strings"
)

type SocialRepository struct {
}

// SocialRepo ...
func SocialRepo() providers.SocialRepository {
	return &SocialRepository{}
}

// ExchangeGoogleOauth2Code ...
func (repo *SocialRepository) ExchangeGoogleOauth2Code(ctx context.Context, data map[string]interface{}) (*models.Traits, error) {

	if data["code"] == nil || data["code"] == "" {
		return nil, models.ErrProviderToken
	}

	params := fmt.Sprintf("?code=%v&client_id=%v&client_secret=%v&redirect_uri=%v&grant_type=authorization_code",
		data["code"],
		env.GetString("apis.providers.google.credential.client_id"),
		env.GetString("apis.providers.google.credential.client_secret"),
		env.GetString("apis.providers.google.credential.redirect_url"))

	ExchangeUrl := env.GetString("apis.providers.google.exchange_code_url.url") + params

	exchangeMethod := env.GetString("apis.providers.google.exchange_code_url.method")

	res, err := helpers.SendHttpRequest(ExchangeUrl, exchangeMethod, nil, nil)

	logger.Debug(res)
	if err != nil {
		return nil, models.ErrProviderConnection
	}

	info, ok := res.(map[string]interface{})
	if !ok {
		return nil, models.ErrProviderConnection
	}

	if info["id_token"] == nil || info["id_token"] == "" {
		return nil, models.ErrProviderConnection
	}

	payload, err := helpers.ExtractClaimsFromTokenString(fmt.Sprintf("%v", info["id_token"]))

	logger.Debug(payload)
	if err != nil {
		return nil, models.ErrProviderConnection
	}

	//traits := data["user"].(map[string]interface{})
	traits := &models.Traits{
		Email:         fmt.Sprintf("%v", payload["email"]),
		EmailVerified: fmt.Sprintf("%v", payload["email_verified"]),
		DisplayName:   fmt.Sprintf("%v", payload["name"]),
		Name: &models.Name{
			First: fmt.Sprintf("%v", payload["given_name"]),
			Last:  fmt.Sprintf("%v", payload["family_name"]),
		},
		Avatar:   fmt.Sprintf("%v", payload["picture"]),
		Username: strings.ToLower(helpers.GenerateRandomString(8)),
	}
	logger.Debug(traits)
	return traits, nil
}

// ExchangeAppleOauth2Code ...
func (repo *SocialRepository) ExchangeAppleOauth2Code(ctx context.Context, data map[string]interface{}) (*models.Traits, error) {

	if data["code"] == nil || data["code"] == "" {
		return nil, models.ErrProviderToken
	}

	params := fmt.Sprintf("?code=%v&client_id=%v&client_secret=%v&redirect_uri=%v&grant_type=authorization_code",
		data["code"],
		env.GetString("apis.providers.google.credential.client_id"),
		env.GetString("apis.providers.google.credential.client_secret"),
		env.GetString("apis.providers.google.credential.redirect_url"))

	ExchangeUrl := env.GetString("apis.providers.google.exchange_code_url.url") + params

	exchangeMethod := env.GetString("apis.providers.google.exchange_code_url.method")

	res, err := helpers.SendHttpRequest(ExchangeUrl, exchangeMethod, nil, nil)

	logger.Debug(res)
	if err != nil {
		return nil, models.ErrProviderConnection
	}

	info, ok := res.(map[string]interface{})
	if !ok {
		return nil, models.ErrProviderConnection
	}

	if info["id_token"] == nil || info["id_token"] == "" {
		return nil, models.ErrProviderConnection
	}

	payload, err := helpers.ExtractClaimsFromTokenString(fmt.Sprintf("%v", info["id_token"]))

	logger.Debug(payload)
	if err != nil {
		return nil, models.ErrProviderConnection
	}

	//traits := data["user"].(map[string]interface{})
	traits := &models.Traits{
		Email:         fmt.Sprintf("%v", payload["email"]),
		EmailVerified: fmt.Sprintf("%v", payload["email_verified"]),
		DisplayName:   fmt.Sprintf("%v", payload["name"]),
		Name: &models.Name{
			First: fmt.Sprintf("%v", payload["given_name"]),
			Last:  fmt.Sprintf("%v", payload["family_name"]),
		},
		Avatar: fmt.Sprintf("%v", payload["picture"]),
	}
	logger.Debug(traits)
	return traits, nil
}

// VerifyGoogle ...
func (repo *SocialRepository) VerifyGoogle(ctx context.Context, data map[string]interface{}) (*models.Traits, error) {

	if data["idToken"] == nil || data["idToken"] == "" {
		return nil, models.ErrProviderToken
	}

	if data["user"] == nil || data["user"] == "" {
		return nil, models.ErrProviderToken
	}

	verifyUrl := env.GetString("apis.providers.google.verify_url.url") + "?id_token=" + fmt.Sprintf("%v", data["idToken"])
	verifyMethod := env.GetString("apis.providers.google.verify_url.method")

	res, err := helpers.SendHttpRequest(verifyUrl, verifyMethod, nil, nil)

	if err != nil {
		return nil, models.ErrProviderConnection
	}

	info, ok := res.(map[string]interface{})
	if !ok {
		return nil, models.ErrProviderConnection
	}

	if info["error"] != nil || info["error"] != "" {
		return nil, models.ErrProviderConnection
	}

	if info["content"] == nil || info["content"] != "" {
		return nil, models.ErrProviderConnection
	}

	payload, err := helpers.ExtractClaimsFromTokenString(fmt.Sprintf("%v", data["id_token"]))

	logger.Error(payload)
	if err != nil {
		return nil, models.ErrProviderConnection
	}

	//traits := data["user"].(map[string]interface{})
	traits := &models.Traits{
		Email:         fmt.Sprintf("%v", payload["email"]),
		EmailVerified: fmt.Sprintf("%v", payload["email_verified"]),
		DisplayName:   fmt.Sprintf("%v", payload["name"]),
		Name: &models.Name{
			First: fmt.Sprintf("%v", payload["given_name"]),
			Last:  fmt.Sprintf("%v", payload["family_name"]),
		},
		Avatar:   fmt.Sprintf("%v", payload["picture"]),
		Username: strings.ToLower(helpers.GenerateRandomString(8)),
	}

	return traits, nil
}

// VerifyApple ...
func (repo *SocialRepository) VerifyApple(ctx context.Context, data map[string]interface{}) (*models.Traits, error) {

	if data["identityToken"] == nil || data["identityToken"] == "" {
		return nil, models.ErrProviderToken
	}

	secret, _ := apple.GenerateClientSecret(
		env.GetString("apis.providers.apple.credential.key_file_path"),
		env.GetString("apis.providers.apple.credential.team_id"),
		env.GetString("apis.providers.apple.credential.client_id"),
		env.GetString("apis.providers.apple.credential.key_file_id"))

	// Generate a new validation client
	client := apple.New()
	vReq := apple.AppValidationTokenRequest{
		ClientID:     env.GetString("apis.providers.apple.credential.client_id"),
		ClientSecret: secret,
		Code:         fmt.Sprintf("%v", data["authorizationCode"]),
	}
	var resp apple.ValidationResponse

	// Do the verification
	err := client.VerifyAppToken(context.Background(), vReq, &resp)

	if err != nil {
		return nil, models.ErrProviderConnection
	}
	_, err = apple.GetUniqueID(resp.IDToken)

	if err != nil {
		return nil, models.ErrProviderConnection
	}

	claim, err := apple.GetClaims(resp.IDToken)
	if err != nil {
		return nil, models.ErrProviderConnection
	}
	email := (*claim)["email"]

	traits := &models.Traits{
		Email:         fmt.Sprintf("%v", email),
		EmailVerified: "true",
		Username:      strings.ToLower(helpers.GenerateRandomString(8)),
		DisplayName:   env.GetString("identity.default_display_name"),
	}

	return traits, nil
}
