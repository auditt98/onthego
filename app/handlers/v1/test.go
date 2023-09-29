package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/auditt98/onthego/utils"
	"github.com/gin-gonic/gin"
)

type TestHandler struct{}

func (ctrl TestHandler) TestPublic(c *gin.Context) {
	fmt.Println("~~~~~" + os.Getenv("HELLO_WORLD"))
	utils.SuccessResponse(c, os.Getenv("HELLO_WORLD"))
	return
}

func (ctrl TestHandler) TestPrivate(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
	return
}
