package zitadel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/auditt98/onthego/types"
	"github.com/go-resty/resty/v2"
)

type SecretAPIResponse struct {
	ClientId     string `json:"ClientId"`
	ClientSecret string `json:"ClientSecret"`
	AppId        string `json:"AppId"`
}

func AddUserToOrg(jwt, userId string, roles []string, orgId string) (bool, error) {
	var err types.ZitadelError
	type AddUserToOrgRequest struct {
		UserId string   `json:"userId"`
		Roles  []string `json:"roles"`
	}

	var addUserToOrgRequest AddUserToOrgRequest
	addUserToOrgRequest.UserId = userId
	addUserToOrgRequest.Roles = roles

	client := resty.New()
	request := client.R().
		ForceContentType("application/json").
		SetBody(addUserToOrgRequest).
		SetAuthToken(jwt).
		SetError(&err)

	if orgId != "" {
		request.SetHeader("x-zitadel-orgid", orgId)
	}

	_, e := request.Post(os.Getenv("ZITADEL_DOMAIN") + "/management/v1/orgs/me/members")

	if e != nil {
		fmt.Println("Error adding user to org: ", e.Error())
		return false, e
	}
	if err.Code != 0 || err.Message != "" {
		fmt.Println("Error adding user to org: ", err.Message)
		return false, fmt.Errorf(err.Message)
	}
	return true, nil
}

func CreateAction(jwt, orgId, name, script string, allowedToFail bool) (string, error) {
	var err types.ZitadelError
	type CreateActionRequest struct {
		Name          string `json:"name"`
		Script        string `json:"script"`
		AllowedToFail bool   `json:"allowedToFail"`
	}

	type CreateActionResponse struct {
		Id string `json:"id"`
	}
	var createActionResponse CreateActionResponse
	var createActionRequest CreateActionRequest
	createActionRequest.Name = name
	createActionRequest.Script = script
	createActionRequest.AllowedToFail = allowedToFail

	client := resty.New()
	request := client.R().
		ForceContentType("application/json").
		SetBody(createActionRequest).
		SetAuthToken(jwt).
		SetResult(&createActionResponse).
		SetError(&err)

	if orgId != "" {
		request.SetHeader("x-zitadel-orgid", orgId)
	}

	_, e := request.Post(os.Getenv("ZITADEL_DOMAIN") + "/management/v1/actions")

	if e != nil {
		fmt.Println("Error adding action: ", e.Error())
		return "", e
	}
	if err.Code != 0 || err.Message != "" {
		fmt.Println("Error adding action: ", err.Message)
		return "", fmt.Errorf(err.Message)
	}
	return createActionResponse.Id, nil
}

func ReadDefaultClientID() string {
	str, err := ioutil.ReadFile(os.Getenv("DEFAULT_CLIENT_ID_PATH"))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(str)
}

func ReadDefaultAPISecret() *SecretAPIResponse {
	// read secret from file

	var apiResponse SecretAPIResponse
	jsonData, err := ioutil.ReadFile(os.Getenv("DEFAULT_API_SECRET_PATH"))

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	err = json.Unmarshal(jsonData, &apiResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return &apiResponse
}

func AddDefaultUserGrantAction(jwt, orgId, projectId string) string {
	apiResponse := ReadDefaultAPISecret()
	res, _ := CreateAction(jwt, orgId, "addGrant", "function addGrant(ctx, api) {let http = require('zitadel/http'); let logger = require('zitadel/log'); logger.log('ctx', ctx.v1); api.userGrants.push({projectID: '"+projectId+"',roles: ['USER']}); var user = http.fetch('"+os.Getenv("API_DOMAIN")+"/api/public/idp/import', {method: 'POST',body: { 'id': ctx.v1.getUser().id, 'clientId': '"+apiResponse.ClientId+"', 'secret': '"+apiResponse.ClientSecret+"'}}).json();logger.log(user.id); }", true)
	return res

	// http.fetch('"+os.Getenv("API_DOMAIN")+"/api/v1/test', {method: 'POST',body: {'id': ctx.v1.getUser().id}, 'clientId': '"+apiResponse.ClientId+"', 'secret': '"+apiResponse.ClientSecret+"' }).json();
}

func SetTriggerAction(jwt, orgId, flowType, triggerType string, actionIds []string) bool {
	var err types.ZitadelError
	type SetTriggerActionRequest struct {
		ActionIds []string `json:"actionIds"`
	}

	var setTriggerActionRequest SetTriggerActionRequest
	setTriggerActionRequest.ActionIds = actionIds

	client := resty.New()
	request := client.R().
		ForceContentType("application/json").
		SetBody(setTriggerActionRequest).
		SetAuthToken(jwt).
		SetError(&err)

	if orgId != "" {
		request.SetHeader("x-zitadel-orgid", orgId)
	}

	_, e := request.Post(os.Getenv("ZITADEL_DOMAIN") + "/management/v1/flows/" + flowType + "/trigger/" + triggerType)

	if e != nil {
		fmt.Println("Error setting trigger action: ", e.Error())
		return false
	}
	if err.Code != 0 || err.Message != "" {
		fmt.Println("Error setting trigger action: ", err.Message)
		return false
	}
	return true
}
