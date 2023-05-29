package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"raptor/app/providers"
	env "raptor/config"
	"raptor/logger"
	"raptor/pkg/helpers"

	"raptor/models"
)

type KratosRepository struct{}

// KratosRepo ...
func KratosRepo() providers.KratosRepository {
	return &KratosRepository{}
}

// CreateKratosLoginFLow ...
func (repo *KratosRepository) CreateKratosLoginFLow(ctx context.Context) (flow string, err error) {
	url := env.GetString("apis.kratos.public.create_login_flow.url")
	method := env.GetString("apis.kratos.public.create_login_flow.method")
	//send request to kratos to create a login flow
	response, err := helpers.SendHttpRequest(url, method, nil, nil)
	if err != nil {
		return "", models.ErrCreateKratosLoginFlow
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok || responseMap["id"] == nil {
		logger.Error(models.ErrCreateKratosLoginFlow.Error())
		return "", models.ErrCreateKratosLoginFlow
	}
	//get flow id from response
	flow = fmt.Sprintf("%v", responseMap["id"])
	return flow, nil
}

// SubmitKratosLoginFlow ...
func (repo *KratosRepository) SubmitKratosLoginFlow(ctx context.Context, flow string, req *models.SubmitKratosLoginRequest) (interface{}, error) {
	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.public.submit_login_flow.url") + "?flow=" + flow
	method := env.GetString("apis.kratos.public.submit_login_flow.method")
	request, err := json.Marshal(req)
	if err != nil {
		logger.Error(models.ErrSubmitLoginFlow.Error())
		return nil, models.ErrSubmitLoginFlow
	}
	response, err := helpers.SendHttpRequest(url, method, request, nil)
	if err != nil {
		return nil, models.ErrSubmitLoginFlow
	}
	return response, nil
}

// CreateKratosRegisterFLow ...
func (repo *KratosRepository) CreateKratosRegisterFLow(ctx context.Context) (flow string, err error) {
	url := env.GetString("apis.kratos.public.create_register_flow.url")
	method := env.GetString("apis.kratos.public.create_register_flow.method")
	//send request to kratos to create a login flow
	response, err := helpers.SendHttpRequest(url, method, nil, nil)
	if err != nil {
		return "", err
	}
	responseMap, ok := response.(map[string]interface{})

	if !ok || responseMap["id"] == nil {
		logger.Error(models.ErrCreateKratosLoginFlow.Error())
		return "", models.ErrCreateKratosLoginFlow
	}
	//get flow id from response
	flow = fmt.Sprintf("%v", responseMap["id"])
	return flow, nil
}

// SubmitKratosRegisterFLow ...
func (repo *KratosRepository) SubmitKratosRegisterFLow(ctx context.Context, flow string, req *models.SubmitKratosRegisterRequest) (interface{}, error) {
	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.public.submit_register_flow.url") + "?flow=" + flow
	method := env.GetString("apis.kratos.public.submit_register_flow.method")
	request, err := json.Marshal(req)
	if err != nil {
		return nil, models.ErrSubmitLoginFlow
	}
	response, err := helpers.SendHttpRequest(url, method, request, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return response, nil
}

// CheckIdentityExistence ...
func (repo *KratosRepository) CheckIdentityExistence(ctx context.Context, identity string, identityType string) (bool, error) {
	query := `query ($jsonFilter: jsonb) {
			identities(where: {traits: {_contains: $jsonFilter}}, limit: 1) {
			id
			traits
			}
		}`
	traits := make(map[string]interface{})
	vars := make(map[string]interface{})
	if identityType == string(models.Email) {
		traits["email"] = identity
	}
	if identityType == string(models.Sms) {
		traits["phone_number"] = identity
	}
	if identityType == string(models.Username) {
		traits["username"] = identity
	}

	vars["jsonFilter"] = traits
	res, err := helpers.SendGraphqlRequest(env.GetString("apis.kratos.graphql.search_identity"), query, vars)
	if err != nil {
		return false, models.ErrConnectToGraphql
	}

	resMap, ok := res.(map[string]interface{})
	if !ok {
		logger.Error(models.ErrConnectToGraphql.Error())
		return false, models.ErrConnectToGraphql
	}

	identities, ok := resMap["identities"].([]interface{})
	if !ok {
		logger.Error(models.ErrConnectToGraphql.Error())
		return false, models.ErrConnectToGraphql
	}

	if len(identities) > 0 {
		return true, nil
	}
	//fmt.Println(res)
	return false, nil
}

// UpdateIdentity ...
func (repo *KratosRepository) UpdateIdentity(ctx context.Context, identityID string, identity interface{}) (interface{}, error) {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.admin.update_identity") + "/" + identityID
	method := "PUT"

	request, err := json.Marshal(identity)
	if err != nil {
		logger.Error(models.ErrSubmitLoginFlow.Error())
		return nil, models.ErrSubmitLoginFlow
	}
	res, err := helpers.SendHttpRequest(url, method, request, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetIdentity ...
func (repo *KratosRepository) GetIdentity(ctx context.Context, claims map[string]interface{}) (interface{}, error) {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.admin.get_identity") + "/" + fmt.Sprintf("%v", claims["sub"])
	method := "GET"

	res, err := helpers.SendHttpRequest(url, method, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ActiveIdentity ...
func (repo *KratosRepository) ActiveIdentity(ctx context.Context, jwtClaims map[string]interface{}, identityFormat string) (interface{}, error) {

	//create
	activate := make(map[string]interface{})
	activate["status"] = models.Active

	resp, err := repo.GetIdentity(ctx, jwtClaims)
	if err != nil {
		return nil, err
	}

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.admin.update_identity") + "/" + fmt.Sprintf("%v", jwtClaims["sub"])
	method := "PUT"

	info := resp.(map[string]interface{})
	info["credentials"] = nil

	if info["metadata_public"] != nil {
		metaData, ok := info["metadata_public"].(map[string]interface{})
		if !ok {
			logger.Error(models.ErrIdentityFormat.Error())
			return nil, models.ErrIdentityFormat
		}
		metaData["status"] = models.Active
	} else {
		info["metadata_public"] = activate
	}
	traits, ok := info["traits"].(map[string]interface{})
	if !ok {
		logger.Error(models.ErrIdentityFormat.Error())
		return nil, models.ErrIdentityFormat
	}
	if identityFormat == string(models.Email) {
		traits["email_verified"] = "true"
	}

	if identityFormat == string(models.Sms) {
		traits["phone_number_verified"] = "true"

	}
	request, err := json.Marshal(resp)
	if err != nil {
		logger.Error(models.ErrIdentityFormat.Error())
		return nil, models.ErrIdentityFormat
	}

	res, err := helpers.SendHttpRequest(url, method, request, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
