package routes

import (
	hv1 "github.com/auditt98/onthego/handlers/v1"
	"github.com/gin-gonic/gin"
)

func V1Router(v1 gin.IRoutes) {
	user := hv1.UserHandlerV1{}
	v1.POST("/idp/import", user.AddUserFromIdP)

	test := hv1.TestHandler{}
	v1.GET("/test/public", test.TestPublic)
	v1.GET("/test/private", test.TestPrivate)
}
