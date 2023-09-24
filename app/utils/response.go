package utils

import (
	"net/http"

	"github.com/auditt98/onthego/types"
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, types.Response{Code: httpStatus, Message: message, Details: data})
	c.Abort()
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, types.Response{Code: http.StatusOK, Message: "", Details: data})
}
