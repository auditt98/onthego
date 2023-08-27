package middlewares

import (
	"strings"

	"github.com/auditt98/onthego/types"
	"github.com/auditt98/onthego/zitadel"
	"github.com/gin-gonic/gin"
)

func TokenIntrospectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, types.ZitadelError{Code: 401, Message: "No Authorization header found"})
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, types.ZitadelError{Code: 401, Message: "Missing or invalid Authorization header"})
			return
		}

		token := authParts[1]
		introspectionResult, error := zitadel.IntrospectToken(token)
		if error != nil {
			c.AbortWithStatusJSON(401, types.ZitadelError{Code: 401, Message: error.Error()})
			return
		}
		if introspectionResult.Active == false {
			c.AbortWithStatusJSON(401, types.ZitadelError{Code: 401, Message: error.Error()})
			return
		}
		c.Set("introspectionResult", introspectionResult)
		c.Next()
	}
}
