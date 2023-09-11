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

type LikeHandlerV1 struct{}

func (ctrl LikeHandlerV1) Like(c *gin.Context) {
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))

	likeValidator := validatorsV1.LikeValidator{}
	if err := c.ShouldBindJSON(&likeValidator); err != nil {
		c.JSON(http.StatusNotAcceptable, types.Error{Code: 406, Message: "Invalid request data"})
		return
	}
	album := models.Album{}
	photo := models.Photo{}
	like := models.Like{}

	if likeValidator.AlbumID != 0 {
		var albumSearchParams = db.SearchParams{
			Filters: map[string]any{
				"id": likeValidator.AlbumID,
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
		var likeSearchParams = db.SearchParams{
			Filters: map[string]any{
				"liker_id":   introspection.Sub,
				"liker_type": "User",
				"album_id":   album.ID,
			},
		}
		db.QueryOne(&likeSearchParams, nil, &like)
		if like.ID != 0 {
			album.Likes = []*models.Like{}
			for _, albumLike := range album.Likes {
				if albumLike.ID != like.ID {
					album.Likes = append(album.Likes, albumLike)
				}
			}
			album.LikesCount = len(album.Likes)
			db.DB.Save(&album)
			db.DB.Delete(&like)
			c.JSON(http.StatusOK, types.SuccessResponse{Data: album})
			return
		}
		like = models.Like{
			LikerID:   introspection.Sub,
			LikerType: "User",
			AlbumID:   album.ID,
		}
		album.Likes = append(album.Likes, &like)
		album.LikesCount = len(album.Likes)
		db.DB.Save(&album)
		c.JSON(http.StatusOK, types.SuccessResponse{Data: album})
		return
	}
	if likeValidator.PhotoID != 0 {
		//check if photo exists
		var photoSearchParams = db.SearchParams{
			Filters: map[string]any{
				"id": likeValidator.PhotoID,
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
		var likeSearchParams = db.SearchParams{
			Filters: map[string]any{
				"liker_id":   introspection.Sub,
				"liker_type": "User",
				"photo_id":   photo.ID,
			},
		}
		db.QueryOne(&likeSearchParams, nil, &like)
		if like.ID != 0 {
			photo.Likes = []*models.Like{}
			for _, photoLike := range photo.Likes {
				if photoLike.ID != like.ID {
					photo.Likes = append(photo.Likes, photoLike)
				}
			}
			photo.LikesCount = len(photo.Likes)
			db.DB.Save(&photo)
			db.DB.Delete(&like)
			c.JSON(http.StatusOK, types.SuccessResponse{Data: photo})
			return
		}
		like = models.Like{
			LikerID:   introspection.Sub,
			LikerType: "User",
			PhotoID:   photo.ID,
		}
		photo.Likes = append(photo.Likes, &like)
		photo.LikesCount = len(photo.Likes)
		db.DB.Save(&photo)
		c.JSON(http.StatusOK, types.SuccessResponse{Data: photo})
		return
	}
	return
}
