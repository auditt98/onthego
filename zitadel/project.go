package zitadel

import (
	"fmt"
	"os"

	"github.com/auditt98/onthego/types"
	"github.com/go-resty/resty/v2"
)

func CreateDefaultProject(jwt, name string, pRoleAssertion, pRoleCheck, hasProjectCheck bool, orgId string) (string, error) {
	var err types.ZitadelError
	type CreateProjectRequest struct {
		Name                 string `json:"name"`
		ProjectRoleAssertion bool   `json:"projectRoleAssertion"`
		ProjectRoleCheck     bool   `json:"projectRoleCheck"`
		HasProjectCheck      bool   `json:"hasProjectCheck"`
	}

	type CreateProjectResponse struct {
		Id string `json:"id"`
	}

	var createProjectRequest CreateProjectRequest
	var createProjectResponse CreateProjectResponse

	createProjectRequest.Name = name
	createProjectRequest.ProjectRoleAssertion = pRoleAssertion
	createProjectRequest.ProjectRoleCheck = pRoleCheck
	createProjectRequest.HasProjectCheck = hasProjectCheck

	client := resty.New()
	request := client.R().
		ForceContentType("application/json").
		SetBody(createProjectRequest).
		SetAuthToken(jwt).
		SetError(&err).
		SetResult(&createProjectResponse)

	if orgId != "" {
		request.SetHeader("x-zitadel-orgid", orgId)
	}

	_, e := request.Post(os.Getenv("ZITADEL_DOMAIN") + "/management/v1/projects")
	if e != nil {
		return "", e
	}
	if err.Code != 0 || err.Message != "" {
		return "", fmt.Errorf(err.Message)
	}
	return createProjectResponse.Id, nil
}
