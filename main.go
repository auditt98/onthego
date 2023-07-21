package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/utils"

	"github.com/gin-contrib/gzip"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
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

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Start the default gin server
	r := gin.Default()

	//Custom form validator
	// binding.Validator = new(forms.DefaultValidator)

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start database
	//Example: db.GetDB() - More info in the models folder
	db.Init()

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis(1)

	params := db.QueryParams{
		Where: db.WhereParams{
			Or: []db.WhereParams{
				{
					Attr: map[string]db.AttributeParams{
						"id": {
							Lt: "5",
						},
					},
				},
			},
		},
	}
	// db.Query("user", params).Find(&results)
	var users []models.User
	db.Query("users", params, &users)
	for _, user := range users {
		fmt.Println("----------")
		fmt.Println(user.ID)
		fmt.Println(user.Email)
		fmt.Println(user.Password)
		fmt.Println(user.Name)

	}
	// generators.LoadAPIVersions(r)

	//generator: Router

	//endgenerator: Router

	// v1 := r.Group("/v1")
	// {
	//generator: Handler
	// }

	// v1 := r.Group("/v1")
	// {
	/*** START USER ***/
	// user := new(handlers.UserController)

	// v1.GET("/test")

	// v1.POST("/user/login", user.Login)
	// v1.POST("/user/register", user.Register)
	// v1.GET("/user/logout", user.Logout)

	// CUSTOM CODE

	// auth := new(handlers.AuthController)

	//Refresh the token when needed to generate new access_token and refresh_token for the user
	// v1.POST("/token/refresh", auth.Refresh)

	/*** START Article ***/
	// article := new(controllers.ArticleController)

	// v1.POST("/article", TokenAuthMiddleware(), article.Create)
	// v1.GET("/articles", TokenAuthMiddleware(), article.All)
	// v1.GET("/article/:id", TokenAuthMiddleware(), article.One)
	// v1.PUT("/article/:id", TokenAuthMiddleware(), article.Update)
	// v1.DELETE("/article/:id", TokenAuthMiddleware(), article.Delete)
	// }

	r.LoadHTMLGlob("./public/html/*")

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03",
			"goVersion":             runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	port := os.Getenv("PORT")

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {

		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer",
			KEY:  "./cert/myCA.key",
		}

		r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		r.Run(":" + port)
	}

}
