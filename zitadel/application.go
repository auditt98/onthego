package zitadel

import (
	"fmt"
	"os"

	"github.com/auditt98/onthego/types"
	"github.com/go-resty/resty/v2"
)

type OIDCResponseType string
type OIDCGrantType string
type OIDCAppType string
type OIDCAuthMethodType string

type CreateOIDCAppRequest struct {
	Name                     string             `json:"name"`
	ResponseTypes            []OIDCResponseType `json:"responseTypes"`
	GrantTypes               []OIDCGrantType    `json:"grantTypes"`
	RedirectURIs             []string           `json:"redirectUris"`
	AppType                  OIDCAppType        `json:"appType"`
	AuthMethodType           OIDCAuthMethodType `json:"authMethodType"`
	DevMode                  bool               `json:"devMode"`
	AccessTokenRoleAssertion bool               `json:"accessTokenRoleAssertion"`
	IDTokenRoleAssertion     bool               `json:"idTokenRoleAssertion"`
	IdTokenUserInfoAssertion bool               `json:"idTokenUserinfoAssertion"`
}

type CreateOIDCAppResponse struct {
	AppId         string `json:"appId"`
	ClientId      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
	NoneCompliant bool   `json:"noneCompliant"`
}

type CreateAPIAppRequest struct {
	Name string `json:"name"`
}

type CreateAPIAppResponse struct {
	AppId        string `json:"appId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func CreateOIDCApp(orgId, projectId, jwt string, oidcAppRequest CreateOIDCAppRequest) (*CreateOIDCAppResponse, error) {
	var err types.ZitadelError
	var response CreateOIDCAppResponse
	client := resty.New()
	request := client.R().
		ForceContentType("application/json").
		SetBody(oidcAppRequest).
		SetAuthToken(jwt).
		SetResult(&response).
		SetError(&err)

	_, e := request.Post(os.Getenv("ZITADEL_DOMAIN") + "/management/v1" + "/projects/" + projectId + "/apps/oidc")

	if e != nil {
		return nil, e
	}
	if err.Code != 0 || err.Message != "" {
		return nil, fmt.Errorf(err.Message)
	}
	return &response, nil
}

func CreateAPIApp(jwt, projectId, name string) (*CreateAPIAppResponse, error) {
	var err types.ZitadelError
	var response CreateAPIAppResponse
	var apiAppRequest CreateAPIAppRequest
	apiAppRequest.Name = name

	client := resty.New()
	request := client.R().
		ForceContentType("application/json").
		SetBody(apiAppRequest).
		SetAuthToken(jwt).
		SetResult(&response).
		SetError(&err)

	_, e := request.Post(os.Getenv("ZITADEL_DOMAIN") + "/management/v1" + "/projects/" + projectId + "/apps/api")

	if e != nil {
		fmt.Println("Error creating default api app ", e.Error())
		return nil, e
	}
	if err.Code != 0 || err.Message != "" {
		fmt.Println("Error creating default api app ", err.Message)
		return nil, fmt.Errorf(err.Message)
	}
	return &response, nil
}
