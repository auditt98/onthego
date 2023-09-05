package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/auditt98/onthego/db"
	hv1 "github.com/auditt98/onthego/handlers/v1"
	hv2 "github.com/auditt98/onthego/handlers/v2"
	"github.com/auditt98/onthego/middlewares"
	"github.com/auditt98/onthego/utils"
	"github.com/auditt98/onthego/zitadel"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkgin "github.com/rookie-ninja/rk-gin/v2/boot"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// auth.TokenValid(c)
		c.Next()
	}
}

func LoadEnv() error {
	args := os.Args
	var env string
	//get env from cmd line argument env=[dev/staging/prod]
	//else get from env variable ONTHEGO_ENV
	//else default to dev
	for _, arg := range args {
		if arg[:4] == "env=" {
			env = arg[4:]
			break
		}
	}

	if env == "" {
		env = os.Getenv("ONTHEGO_ENV")
	}
	if env == "" {
		env = "dev"
	}
	//load .env.[env].local, .env.local, .env.[env], .env
	godotenv.Load(".env." + env + ".local")
	godotenv.Load(".env.local")
	godotenv.Load(".env." + env)
	godotenv.Load()
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	return nil
}

func InitZitadel() {
	k, _ := zitadel.GenerateJWTFromKeyFile()
	userId, err := zitadel.CreateDefaultHumanUser(k)
	if err != nil {
		return
	}
	ok, _ := zitadel.AddUserToOrg(k, userId, []string{"ORG_OWNER"}, "")
	if !ok {
		return
	}

	ok, _ = zitadel.AddUserToIAM(k, userId)
	if !ok {
		return
	}

	defaultProjectId, e := zitadel.CreateDefaultProject(k, "OnTheGo", true, false, true, "")
	if e != nil {
		return
	}
	fmt.Println("Default ProjectId: ", defaultProjectId)
	e = zitadel.BulkAddRoleToProject(k, defaultProjectId, []zitadel.RoleRequest{
		{
			Key:         "ADMIN",
			DisplayName: "Admin",
		},
		{
			Key:         "USER",
			DisplayName: "User",
		},
		{
			Key:         "MODERATOR",
			DisplayName: "Moderator",
		},
	})
	if e != nil {
		return
	}
	var defaultAppRequest = zitadel.CreateOIDCAppRequest{
		Name:                     "OnTheWall",
		DevMode:                  true,
		RedirectURIs:             []string{"http://localhost:3000/api/auth/callback/zitadel"},
		ResponseTypes:            []zitadel.OIDCResponseType{"OIDC_RESPONSE_TYPE_CODE"},
		GrantTypes:               []zitadel.OIDCGrantType{"OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_REFRESH_TOKEN"},
		AppType:                  zitadel.OIDCAppType("OIDC_APP_TYPE_WEB"),
		AuthMethodType:           zitadel.OIDCAuthMethodType("OIDC_AUTH_METHOD_TYPE_NONE"),
		AccessTokenRoleAssertion: true,
		IDTokenRoleAssertion:     true,
		IdTokenUserInfoAssertion: true,
	}

	createAppResponse, e := zitadel.CreateOIDCApp("", defaultProjectId, k, defaultAppRequest)
	if e != nil {
		return
	}
	fmt.Println("Default App ClientId: ", createAppResponse.ClientId)

	createAPIResponse, e := zitadel.CreateAPIApp(k, defaultProjectId, "OnTheWall_API")
	if e != nil {
		return
	}
	//write create api response to file
	fmt.Println("Default API ClientId: ", createAPIResponse.ClientId)
	fmt.Println("Default API ClientSecret", createAPIResponse.ClientSecret)
	fmt.Println("Default API AppId", createAPIResponse.AppId)
	jsonData, err := json.Marshal(createAPIResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = ioutil.WriteFile("./machinekey/default_api_secret.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	actionId := zitadel.AddDefaultUserGrantAction(k, "", defaultProjectId)
	if actionId == "" {
		return
	}
	fmt.Println("Default ActionId: ", actionId)
	setTriggerActionResult := zitadel.SetTriggerAction(k, "", "1", "3", []string{actionId})
	setTriggerActionResult = zitadel.SetTriggerAction(k, "", "3", "3", []string{actionId})
	if !setTriggerActionResult {
		return
	}
}

func main() {
	LoadEnv()
	_, err := utils.LoadConf()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	db.Init()
	db.InitRedis(1)
	boot := rkboot.NewBoot()
	entry := rkgin.GetGinEntry("ginboilerplate")

	InitZitadel()
	v1 := entry.Router.Group("/api/v1")
	{
		user := hv1.UserHandlerV1{}
		v1.POST("/idp/import", user.AddUserFromIdP)
		v1.POST("/test", middlewares.TokenIntrospectionMiddleware(), user.Test)
		v1.POST("/test2", user.TestPublic)

		album := hv1.AlbumHandlerV1{}
		v1.POST("/albums/search", middlewares.TokenIntrospectionMiddleware(), album.Search)
		v1.POST("/albums", middlewares.TokenIntrospectionMiddleware(), album.CreateAlbum)
		v1.POST("/albums/:album_id/users", middlewares.TokenIntrospectionMiddleware(), album.AddUserToAlbum)
		v1.DELETE("/albums/:album_id/users/:user_id", middlewares.TokenIntrospectionMiddleware(), album.RemoveUserFromAlbum)
	}
	v2 := entry.Router.Group("/api/v2")
	{
		article := hv2.ArticleHandlerV2{}
		v2.GET("/test", article.Get)
	}

	entry.Router.LoadHTMLGlob("./public/html/*")
	// logger := rkentry.GlobalAppCtx.GetLoggerEntry("my-logger")
	boot.Bootstrap(context.TODO())
	boot.WaitForShutdownSig(context.TODO())
}

func writeOK(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK " + time.Now().String()))
}
