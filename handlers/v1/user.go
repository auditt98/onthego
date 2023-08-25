package handlers

import (
	"fmt"
	"net/http"

	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
)

type UserHandlerV1 struct{}

func (ctrl UserHandlerV1) AddUserFromIdP(c *gin.Context) {
	fmt.Println("HElloo ADD USER FROM IDP")
	// userImportValidator := validatorsV1.IdPUserImportValidator{}
	// if c.ShouldBindJSON(&userImportValidator) != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid request data", "form": userImportValidator})
	// 	c.Abort()
	// 	return
	// }
	fmt.Println("Hello")
	c.JSON(200, "Hello")
	return
}

func (ctrl UserHandlerV1) Get(c *gin.Context) {
	c.JSON(200, gin.H{"message": "test"})
	return
}

func (ctrl UserHandlerV1) Update(c *gin.Context) {
	userImportValidator := validatorsV1.IdPUserImportValidator{}
	if c.ShouldBindJSON(&userImportValidator) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid request data", "form": userImportValidator})
		c.Abort()
		return
	}
	fmt.Println(userImportValidator)
	c.JSON(200, userImportValidator)
	return
}

func (ctrl UserHandlerV1) Delete(c *gin.Context) {
	//code here
	return
}
