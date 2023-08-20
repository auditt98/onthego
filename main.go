package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/auditt98/onthego/db"
	hv1 "github.com/auditt98/onthego/handlers/v1"
	hv2 "github.com/auditt98/onthego/handlers/v2"
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
		fmt.Println("Error creating default human user:", err)
		return
	}
	ok, _ := zitadel.AddUserToOrg(k, userId, []string{"ORG_OWNER"}, "")
	if !ok {
		fmt.Println("Error adding user to org")
		return
	}
	defaultProjectId, e := zitadel.CreateDefaultProject(k, "OnTheGo", true, true, true, "")
	if e != nil {
		fmt.Println("Error creating default project:", e)
		return
	}
	fmt.Println("Default project ID: ", defaultProjectId)
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

	//router
	v1 := entry.Router.Group("/api/v1")
	{
		article := hv1.ArticleHandlerV1{}
		v1.GET("/test", article.Get)
		v1.POST("/test", article.Update)
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
