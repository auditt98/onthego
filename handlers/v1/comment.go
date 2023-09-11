package handlers

import (
	"net/http"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	"github.com/auditt98/onthego/utils"
	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
)

type CommentHandlerV1 struct{}

func (ctrl CommentHandlerV1) Comment(c *gin.Context) {
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	commentValidator := validatorsV1.CommentValidator{}
	if err := c.ShouldBindJSON(&commentValidator); err != nil {
		c.JSON(http.StatusNotAcceptable, types.Error{Code: 406, Message: "Invalid request data"})
		return
	}
	album := models.Album{}
	photo := models.Photo{}
	comment := models.Comment{}

	if commentValidator.AlbumID != 0 {
		var albumSearchParams = db.SearchParams{
			Filters: map[string]any{
				"id": commentValidator.AlbumID,
				"users": map[string]any{
					"id": introspection.Sub,
				},
			},
			Populate: []string{"Likes"},
		}
		db.QueryOne(&albumSearchParams, nil, &album)
		if album.ID == 0 {
			c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "Album not found"})
			return
		}
		comment = models.Comment{
			CommenterID:   introspection.Sub,
			CommenterType: "User",
			AlbumID:       album.ID,
			Content:       commentValidator.Content,
		}
		album.Comments = append(album.Comments, &comment)
		db.DB.Save(&album)
		c.JSON(http.StatusOK, types.SuccessResponse{Data: album})
		return
	}
	if commentValidator.PhotoID != 0 {
		var photoSearchParams = db.SearchParams{
			Filters: map[string]any{
				"id": commentValidator.PhotoID,
			},
		}
		db.QueryOne(&photoSearchParams, nil, &photo)
		if photo.ID == 0 {
			c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "Photo not found"})
			return
		}
		albumSearchParams := db.SearchParams{
			Filters: map[string]any{
				"id": photo.AlbumID,
				"users": map[string]any{
					"id": introspection.Sub,
				},
			},
		}
		db.QueryOne(&albumSearchParams, nil, &album)
		if album.ID == 0 {
			c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "You don't have permission to perform this action"})
			return
		}
		comment = models.Comment{
			CommenterID:   introspection.Sub,
			CommenterType: "User",
			PhotoID:       photo.ID,
			Content:       commentValidator.Content,
		}
		photo.Comments = append(photo.Comments, &comment)
		db.DB.Save(&photo)
		c.JSON(http.StatusOK, types.SuccessResponse{Data: photo})
		return
	}
}

func (ctrl CommentHandlerV1) UpdateComment(c *gin.Context) {
	return
}

func (ctrl CommentHandlerV1) DeleteComment(c *gin.Context) {
	return
}
