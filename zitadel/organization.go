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
