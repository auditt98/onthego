package zitadel

import (
	"fmt"
	"os"

	"github.com/auditt98/onthego/types"
	"github.com/go-resty/resty/v2"
)

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

func AddDefaultUserGrantAction(jwt, orgId, projectId string) string {
	res, _ := CreateAction(jwt, orgId, "addGrant", "function addGrant(ctx, api) {api.userGrants.push({projectID: '"+projectId+"',roles: ['USER']});}", true)
	return res
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
