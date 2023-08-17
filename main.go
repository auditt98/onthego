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
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkgin "github.com/rookie-ninja/rk-gin/v2/boot"

	uuid "github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware ...
// CORS (Cross-Origin Resource Sharing)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// RequestIDMiddleware ...
// Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

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
