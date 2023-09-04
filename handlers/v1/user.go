package handlers

import (
	"fmt"
	"net/http"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/auditt98/onthego/zitadel"
	"github.com/gin-gonic/gin"
)

type UserHandlerV1 struct{}

func (ctrl UserHandlerV1) AddUserFromIdP(c *gin.Context) {
	userImportValidator := validatorsV1.IdPUserImportValidator{}
	if c.ShouldBindJSON(&userImportValidator) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid request data", "form": userImportValidator})
		c.Abort()
		return
	}
	apiSecret := zitadel.ReadDefaultAPISecret()
	if userImportValidator.ClientId != apiSecret.ClientId || userImportValidator.Secret != apiSecret.ClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		c.Abort()
		return
	}
	newUser := models.User{}
	newUser.ID = userImportValidator.ID
	instance, err := db.ResolveDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error resolving database"})
		c.Abort()
		return
	}

	result := instance.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error importing user"})
		c.Abort()
		return
	}
	c.JSON(200, "New user imported")
	return
}

func (ctrl UserHandlerV1) Test(c *gin.Context) {
	introspectionResult, ok := c.Get("introspectionResult")
	if !ok {
		c.JSON(500, types.Error{Code: 500, Message: "Introspection result not found"})
	} else {
		c.JSON(200, types.SuccessResponse{Data: introspectionResult})
	}
	return
}

func (ctrl UserHandlerV1) TestPublic(c *gin.Context) {
	queryObject, _ := c.Get("queryObject")
	// albums := []models.Album{}
	// res := db.ToQuery(queryObject, &albums)
	// c.JSON(200, types.SuccessResponse{Data: map[string]interface{}{"query": query, "queryObject": queryObject}})
	// filter2 := map[string]any{
	// 	"id": "895409660164210689",
	// }
	// var albums2 []models.Album
	// str := db.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Where(filter2).Find(&albums2)
	// })
	c.JSON(200, queryObject)
	// c.JSON(200, res)
	// c.JSON(200, albums)
	// c.JSON(200, query)

	return
	// introspectionResult, ok := c.Get("introspectionResult")
	// filters := map[string]any{
	// 	"name": []string{"Album 1", "Album 2"},
	// }
	// filter2 := map[string]any{
	// 	"users": map[string]any{
	// 		"id": "229409946817527811",
	// 	},
	// }
	// var albums []models.Album
	// result := db.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Or(filters).Or(filter2).Find(&albums)
	// })
	// c.JSON(200, result)
	// c.JSON(200, types.SuccessResponse{Data: albums})
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
