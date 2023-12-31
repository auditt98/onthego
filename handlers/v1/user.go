package handlers

import (
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

func (ctrl UserHandlerV1) GetDefaultClientId(c *gin.Context) {
	response := zitadel.ReadDefaultClientID()
	c.JSON(http.StatusOK, types.SuccessResponse{Data: response})
	return
}
