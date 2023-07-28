package handlers

import (
	"net/http"

	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
)

type ArticleHandlerV1 struct{}

//generic validator

// generator: Actions
func (ctrl ArticleHandlerV1) Create(c *gin.Context) {
	//code here
	return
}

// Greeter handler
// @Summary Greeter
// @Id 1
// @Tags Hello
// @version 1.0
// @Param name query string true "name"
// @produce application/json
// @Success 200 {object} string
// @Router /v1/test [get]
func (ctrl ArticleHandlerV1) Get(c *gin.Context) {
	c.JSON(200, gin.H{"message": "test"})
	return
}

func (ctrl ArticleHandlerV1) Update(c *gin.Context) {
	articleValidator := validatorsV1.ArticleValidator{}
	if c.ShouldBindJSON(&articleValidator) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid request data", "form": articleValidator})
		c.Abort()
		return
	}
	c.JSON(200, articleValidator)
	return
}

func (ctrl ArticleHandlerV1) Delete(c *gin.Context) {
	//code here
	return
}

//endgenerator: Actions
