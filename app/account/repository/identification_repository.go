package repository

import (
	"context"
	"raptor/app/account"
	env "raptor/config"
	"raptor/pkg/helpers"
)

type identificationRepository struct{}

// IdentificationRepo ...
func IdentificationRepo() account.IdentificationRepository {
	return &identificationRepository{}
}

// CheckIdentification ...
func (repo *identificationRepository) CheckIdentification(ctx context.Context, jwt string) (bool, error) {

	url := env.GetString("apis.identification.check_status.url")
	method := env.GetString("apis.identification.check_status.method")
	var bearer = "Bearer " + jwt
	//set header for request
	header := make(map[string]string)
	header["Authorization"] = bearer
	//send request to kratos to create a login res
	response, err := helpers.SendHttpRequest(url, method, nil, header)
	if err != nil {
		return false, err
	}
	responseMap := response.(map[string]interface{})
	data := responseMap["data"].(map[string]interface{})
	//get res id from response
	if data["identified"] == false {
		return false, nil
	}
	return true, nil
}
