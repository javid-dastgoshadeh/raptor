package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"raptor/app/account"
	env "raptor/config"
	"raptor/logger"
	"raptor/pkg/helpers"

	"raptor/models"
)

type kratosRepository struct{}

// KratosRepo ...
func KratosRepo() account.KratosRepository {
	return &kratosRepository{}
}

// CreateSettingFlow ...
func (repo *kratosRepository) CreateSettingFlow(ctx context.Context, claims map[string]interface{}) (string, error) {

	url := env.GetString("apis.kratos.public.create_setting_flow.url")
	method := env.GetString("apis.kratos.public.create_setting_flow.method")

	//set header for request
	header := make(map[string]string)
	header["X-Session-Token"] = claims["sid"].(string)
	//send request to kratos to create a login flow
	response, err := helpers.SendHttpRequest(url, method, nil, header)
	if err != nil {
		return "", err
	}
	responseMap := response.(map[string]interface{})
	//get flow id from response
	flow := fmt.Sprintf("%v", responseMap["id"])
	return flow, nil
}

// SubmitSettingFlow ...
func (repo *kratosRepository) SubmitSettingFlow(ctx context.Context, req *models.SubmitKratosSettingRequest, flow string, claims map[string]interface{}) (interface{}, error) {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.public.submit_setting_flow.url") + "?flow=" + flow
	method := env.GetString("apis.kratos.public.submit_setting_flow.method")

	request, err := json.Marshal(req)
	if err != nil {
		return nil, models.ErrSubmitLoginFlow
	}

	//set header for request
	header := make(map[string]string)
	header["X-Session-Token"] = claims["sid"].(string)
	response, err := helpers.SendHttpRequest(url, method, request, header)

	if err != nil {
		return nil, err
	}
	err = helpers.ExportKratosFlowsErr(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// CheckSession ...
func (repo *kratosRepository) CheckSession(ctx context.Context, claims map[string]interface{}) (interface{}, error) {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.public.check_session.url")
	method := env.GetString("apis.kratos.public.check_session.method")

	//set header for request
	header := make(map[string]string)
	header["X-Session-Token"] = claims["sid"].(string)
	response, err := helpers.SendHttpRequest(url, method, nil, header)
	if err != nil {
		return nil, err
	}
	res, ok := response.(map[string]interface{})
	if !ok {
		logger.Error(models.ErrIdentityFormat.Error())
		return nil, models.ErrIdentityFormat
	}
	if res["error"] != nil {
		logger.Error(models.ErrUnauthorized.Error())
		return nil, models.ErrUnauthorized
	}
	return response, nil
}

// DisableSession ...
func (repo *kratosRepository) DisableSession(ctx context.Context, claims map[string]interface{}) error {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.public.disable_session.url")
	method := env.GetString("apis.kratos.public.disable_session.method")

	//set header for request
	header := make(map[string]string)
	header["X-Session-Token"] = claims["sid"].(string)

	req := make(map[string]string)
	req["session_token"] = claims["sid"].(string)
	request, err := json.Marshal(req)
	if err != nil {
		logger.Error(models.ErrSubmitLoginFlow.Error())
		return models.ErrSubmitLoginFlow
	}

	_, err = helpers.SendHttpRequest(url, method, request, header)
	if err != nil {
		return err
	}
	return nil
}

// InactiveIdentity ...
func (repo *kratosRepository) InactiveIdentity(ctx context.Context, claims map[string]interface{}) error {

	resp, err := repo.GetIdentity(ctx, claims)
	fmt.Println(resp)
	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.admin.update_identity") + "/" + fmt.Sprintf("%v", claims["sub"])
	method := "PUT"

	info := resp.(map[string]interface{})
	info["credentials"] = nil
	metaData := info["metadata_public"].(map[string]interface{})
	metaData["status"] = models.Inactive

	request, err := json.Marshal(resp)

	_, err = helpers.SendHttpRequest(url, method, request, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetIdentity ...
func (repo *kratosRepository) GetIdentity(ctx context.Context, claims map[string]interface{}) (interface{}, error) {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.admin.get_identity") + "/" + fmt.Sprintf("%v", claims["sub"])
	method := "GET"

	res, err := helpers.SendHttpRequest(url, method, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateIdentity ...
func (repo *kratosRepository) UpdateIdentity(ctx context.Context, identityID string, identity interface{}) (interface{}, error) {

	//get login submit url and add the flow id to it as a query param
	url := env.GetString("apis.kratos.admin.update_identity") + "/" + identityID
	method := "PUT"

	request, err := json.Marshal(identity)
	if err != nil {
		return nil, models.ErrSubmitLoginFlow
	}
	res, err := helpers.SendHttpRequest(url, method, request, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
