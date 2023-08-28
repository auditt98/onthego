package handlers

import (
	"net/http"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
)

type AlbumHandlerV1 struct{}

func (ctrl AlbumHandlerV1) CreateAlbum(c *gin.Context) {
	albumValidator := validatorsV1.AlbumValidator{}
	if err := c.ShouldBindJSON(&albumValidator); err != nil {
		c.JSON(http.StatusNotAcceptable, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusNotAcceptable,
			Message: err.Error(),
			Details: err,
		}})
		c.Abort()
		return
	}
	newAlbum := models.Album{}
	newAlbum.Name = albumValidator.Name
	db.DB.Create(&newAlbum)

	return
}
