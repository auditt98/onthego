package handlers

import (
	"net/http"

	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
)

type UserHandlerV1 struct{}

func (ctrl UserHandlerV1) AddUserFromIdP(c *gin.Context) {
	//code here
	return
}

func (ctrl UserHandlerV1) Get(c *gin.Context) {
	c.JSON(200, gin.H{"message": "test"})
	return
}

func (ctrl UserHandlerV1) Update(c *gin.Context) {
	articleValidator := validatorsV1.ArticleValidator{}
	if c.ShouldBindJSON(&articleValidator) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid request data", "form": articleValidator})
		c.Abort()
		return
	}
	c.JSON(200, articleValidator)
	return
}

func (ctrl UserHandlerV1) Delete(c *gin.Context) {
	//code here
	return
}
