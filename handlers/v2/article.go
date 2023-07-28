package handlers

import (
	"github.com/gin-gonic/gin"
)

type ArticleHandlerV2 struct{}

// generator: Actions
func (ctrl ArticleHandlerV2) Create(c *gin.Context) {
	//code here
	return
}

func (ctrl ArticleHandlerV2) Get(c *gin.Context) {
	//code here

	c.JSON(200, gin.H{"message": "test"})
	return
}

func (ctrl ArticleHandlerV2) Update(c *gin.Context) {

	//code here
	return
}

func (ctrl ArticleHandlerV2) Delete(c *gin.Context) {
	//code here
	return
}

//endgenerator: Actions
