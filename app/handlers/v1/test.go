package handlers

import (
	"net/http"

	"github.com/auditt98/onthego/utils"
	"github.com/gin-gonic/gin"
)

type TestHandler struct{}

func (ctrl TestHandler) TestPublic(c *gin.Context) {
	utils.SuccessResponse(c, "Hello World")
	return
}

func (ctrl TestHandler) TestPrivate(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
	return
}
