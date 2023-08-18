package utils

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func MakeGet(url string, headers map[string]string, queryParams map[string]string, token string, responseBody interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R().SetQueryParams(queryParams).SetHeaders(headers).SetResult(responseBody)
	if token != "" {
		req.SetAuthToken(token)
	}

	resp, err := req.Get(url)
	return resp, err
}

func MakeFormPost(url string, headers map[string]string, queryParams map[string]string, data map[string]string, token string, responseBody interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeaders(headers).
		SetQueryParams(queryParams).
		SetFormData(data).
		ForceContentType("application/x-www-form-urlencoded").
		SetResult(responseBody)

	if token != "" {
		req.SetAuthToken(token)
	}
	resp, err := req.Post(url)
	fmt.Println("------RESP--------, ", resp.String())
	fmt.Println("------RBD--------, ", &responseBody)
	return resp, err
}

func MakePost(url string, headers map[string]string, queryParams map[string]string, data interface{}, token string, responseBody interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).
		SetQueryParams(queryParams).
		SetBody(data).
		ForceContentType("application/json").
		SetResult(&responseBody)

	if token != "" {
		req.SetAuthToken(token)
	}
	resp, err := req.Post(url)
	return resp, err
}

func MakePut(url string, headers map[string]string, queryParams map[string]string, data interface{}, token string, responseBody interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R().
		SetHeaders(headers).
		SetQueryParams(queryParams).
		SetBody(data).
		ForceContentType("application/json").
		SetResult(&responseBody)

	if token != "" {
		req.SetAuthToken(token)
	}
	resp, err := req.Put(url)
	return resp, err
}

func MakePatch(url string, headers map[string]string, queryParams map[string]string, data interface{}, token string, responseBody interface{}) (*resty.Response, error) {
	client := resty.New()

	// POST JSON string
	// No need to set content type, if you have client level setting
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParams).
		SetHeaders(headers).
		SetBody(data).
		ForceContentType("application/json").
		SetResult(&responseBody)

	if token != "" {
		req.SetAuthToken(token)
	}
	resp, err := req.Patch(url)
	return resp, err
}

func MakeDelete(url string, headers map[string]string, queryParams map[string]string, token string, responseBody interface{}) (*resty.Response, error) {
	client := resty.New()

	// POST JSON string
	// No need to set content type, if you have client level setting
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).
		ForceContentType("application/json").
		SetResult(&responseBody).
		SetQueryParams(queryParams)

	if token != "" {
		req.SetAuthToken(token)
	}
	resp, err := req.Delete(url)
	return resp, err
}
