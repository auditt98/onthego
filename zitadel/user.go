package zitadel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/auditt98/onthego/types"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	Type     string `json:"type"`
	KeyId    string `json:"keyId"`
	Key      string `json:"key"`
	UserId   string `json:"userId"`
	ClientId string `json:"clientId"`
}

type CreateZitadelUserRequest struct {
	Username string                    `json:"username"`
	Profile  CreateZitadelUserProfile  `json:"profile"`
	Email    CreateZitadelUserEmail    `json:"email"`
	Password CreateZitadelUserPassword `json:"password"`
}

type CreateZitadelUserProfile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateZitadelUserEmail struct {
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
}

type CreateZitadelUserPassword struct {
	Password       string `json:"password"`
	ChangeRequired bool   `json:"changeRequired"`
}

func GenerateJWTFromKey(filePath string, useId string) (string, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
		return "", err
	}

	var keyData JWT
	if err := json.Unmarshal(fileContent, &keyData); err != nil {
		log.Fatal("Error parsing JSON:", err)
		return "", err
	}
	var (
		t *jwt.Token
		s string
	)
	parsedKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData.Key))
	claims := jwt.MapClaims{
		"aud": "http://localhost:8080",
		"iat": time.Now().UTC().Unix(),
		"exp": time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	if useId == "user" {
		claims["iss"] = keyData.UserId
		claims["sub"] = keyData.UserId
	} else if useId == "client" {
		claims["iss"] = keyData.ClientId
		claims["sub"] = keyData.ClientId
	}
	t = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = keyData.KeyId
	s, err = t.SignedString(parsedKey)
	if err != nil {
		log.Fatal("Error signing token:", err)
		return "", err
	}
	return s, nil
}

func GenerateIntrospectionJWT() (string, error) {
	jwt, err := GenerateJWTFromKey("./machinekey/default_api_introspection_secret.json", "client")
	if err != nil {
		log.Fatal("Error signing token:", err)
		return "", err
	}
	return jwt, nil
}

func GenerateJWTServiceUser() (string, error) {
	jwt, err := GenerateJWTFromKey("./machinekey/core_service_user_key.json", "user")
	if err != nil {
		log.Fatal("Error signing token:", err)
		return "", err
	}
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Set("scope", "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud")
	data.Set("assertion", jwt)
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	var token TokenResponse
	_, err = resty.New().R().
		SetResult(&token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
			"scope":      "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud",
			"assertion":  jwt,
		}).
		Post("http://localhost:8080/oauth/v2/token")
	if err != nil {
		log.Fatal("Error making POST request:", err)
		return "", err
	}
	return token.AccessToken, nil
}

func GenerateJWTFromKeyFile() (string, error) {
	filePath := "./machinekey/core_service_user_key.json"
	// Read the JSON file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
		return "", err
	}

	// Parse the JSON into a JWTKey struct
	var keyData JWT
	if err := json.Unmarshal(fileContent, &keyData); err != nil {
		log.Fatal("Error parsing JSON:", err)
		return "", err
	}

	var (
		t *jwt.Token
		s string
	)
	parsedKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData.Key))
	t = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": keyData.UserId,
		"sub": keyData.UserId,
		"aud": "http://localhost:8080",
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().Add(time.Hour * 48).UTC().Unix(),
	})
	t.Header["kid"] = keyData.KeyId
	s, err = t.SignedString(parsedKey)
	if err != nil {
		log.Fatal("Error signing token:", err)
		return "", err
	}

	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Set("scope", "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud")
	data.Set("assertion", s)

	type TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	var token TokenResponse
	_, err = resty.New().R().
		SetResult(&token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
			"scope":      "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud",
			"assertion":  s,
		}).
		Post("http://localhost:8080/oauth/v2/token")
	if err != nil {
		log.Fatal("Error making POST request:", err)
		return "", err
	}
	return token.AccessToken, nil
}

func CheckDefaultHumanUserUnique(jwt string) bool {
	url := os.Getenv("ZITADEL_DOMAIN") + "/management/v1/users/_is_unique"
	//if env: ZITADEL_USER_EMAIL and ZITADEL_USERNAME
	type UniqueResponse struct {
		IsUnique bool `json:"isUnique"`
	}

	var uniqueResponse UniqueResponse

	_, err := resty.New().R().
		ForceContentType("application/json").
		SetAuthToken(jwt).
		SetResult(&uniqueResponse).
		SetQueryParams(map[string]string{
			"email": os.Getenv("ZITADEL_USER_EMAIL"),
		}).Get(url)

	if err != nil {
		fmt.Println("Error making POST request:", err)
		return false
	}
	return uniqueResponse.IsUnique
}

func CreateDefaultHumanUser(jwt string) (string, error) {
	//if env: ZITADEL_USER_EMAIL doesnt exist, return error
	userEmailEnv := os.Getenv("ZITADEL_USER_EMAIL")
	usernameEnv := os.Getenv("ZITADEL_USERNAME")
	passwordEnv := os.Getenv("ZITADEL_PASSWORD")

	if userEmailEnv == "" {
		return "", fmt.Errorf("Missing env variable ZITADEL_USER_EMAIL. Default human user will not be created")
	}

	if passwordEnv == "" {
		return "", fmt.Errorf("Missing env variable ZITADEL_PASSWORD. Default human user will not be created")
	}
	type CreateUserResponse struct {
		UserId string `json:"userId"`
	}

	var createUserResponse CreateUserResponse
	var createUserError types.ZitadelError
	//check if user already exists
	if CheckDefaultHumanUserUnique(jwt) == true {
		var userRequest CreateZitadelUserRequest
		if usernameEnv != "" {
			userRequest.Username = usernameEnv
		}
		userRequest.Profile.FirstName = "Core"
		userRequest.Profile.LastName = "Human User"
		userRequest.Email.Email = userEmailEnv
		userRequest.Email.IsVerified = true
		userRequest.Password.Password = passwordEnv
		userRequest.Password.ChangeRequired = true

		_, err := resty.New().R().
			ForceContentType("application/json").
			SetBody(userRequest).
			SetAuthToken(jwt).
			SetResult(&createUserResponse).
			SetError(&createUserError).
			Post(os.Getenv("ZITADEL_DOMAIN") + "/v2alpha/users/human")

		if err != nil {
			fmt.Println("Error creating default human user ", err.Error())
			return "", err
		}

		if (createUserError.Code != 0) || (createUserError.Message != "") {
			fmt.Println("Error creating default human user ", createUserError.Message)
			return "", fmt.Errorf(createUserError.Message)
		}
		fmt.Println("Default human user created successfully, ID: ", createUserResponse.UserId)

	} else {
		fmt.Println("Default human user already exists. Skipping creation...")
		return "", fmt.Errorf("Default human user already exists. Skipping creation...")
	}
	return createUserResponse.UserId, nil
}

func AddUserToIAM(jwt, userId string) (bool, error) {
	type AddUserToIAMRequest struct {
		UserId string   `json:"userId"`
		Roles  []string `json:"roles"`
	}

	var err types.ZitadelError
	var addUserToIAMRequest AddUserToIAMRequest
	addUserToIAMRequest.UserId = userId
	addUserToIAMRequest.Roles = []string{"IAM_OWNER"}

	_, e := resty.New().R().
		ForceContentType("application/json").
		SetBody(addUserToIAMRequest).
		SetAuthToken(jwt).
		SetError(&err).
		Post(os.Getenv("ZITADEL_DOMAIN") + "/admin/v1/members")

	if e != nil {
		fmt.Println("Error creating default human user ", e.Error())
		return false, e
	}

	if (err.Code != 0) || (err.Message != "") {
		fmt.Println("Error creating default human user ", err.Message)
		return false, fmt.Errorf(err.Message)
	}
	return true, nil
}

func VerifySecret(clientId, secret string) bool {
	type SecretAPIResponse struct {
		ClientId     string `json:"ClientId"`
		ClientSecret string `json:"ClientSecret"`
		AppId        string `json:"AppId"`
	}

	var apiResponse SecretAPIResponse

	jsonData, err := ioutil.ReadFile("./machinekey/default_api_secret.json")
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	err = json.Unmarshal(jsonData, &apiResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	if apiResponse.ClientId == clientId && apiResponse.ClientSecret == secret {
		return true
	} else {
		return false
	}
}

func generateBasicAuthHeader(clientID, clientSecret string) string {
	encodedClientID := url.QueryEscape(clientID)
	encodedClientSecret := url.QueryEscape(clientSecret)

	credentials := encodedClientID + ":" + encodedClientSecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	return "Basic " + encodedCredentials
}

func IntrospectToken(token string) (*types.IntrospectionResult, error) {
	jwt, _ := GenerateIntrospectionJWT()

	var introspectResponse types.IntrospectionResult
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"token":                 token,
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"client_assertion":      jwt,
		}).
		SetResult(&introspectResponse).
		Post(os.Getenv("ZITADEL_DOMAIN") + "/oauth/v2/introspect")
	if err != nil {
		fmt.Println("Invalid client token", err)
		return nil, err
	}
	if introspectResponse.Active == false {
		return nil, fmt.Errorf("Invalid token")
	}
	return &introspectResponse, nil
}
