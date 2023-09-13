package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	"github.com/auditt98/onthego/utils"
	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
)

type PhotoHandlerV1 struct{}

func (ctrl PhotoHandlerV1) Search(c *gin.Context) {
	//get photos
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	userID := introspection.Sub
	currentUserFilter := map[string]any{
		"user_id": userID,
	}
	photos := []models.Photo{}

	searchParams := db.SearchParams{}
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		db.DB.Where(currentUserFilter).Find(&photos)
	}
	var count int64
	db.Query(&searchParams, currentUserFilter, &photos, &count)

	if os.Getenv("UPLOAD_DRIVER") == "local" {
		photoWithPresignedUrl := []map[string]interface{}{}
		var scheme string
		if c.Request.Header.Get("X-Forwarded-Proto") != "" {
			scheme = c.Request.Header.Get("X-Forwarded-Proto")
		} else if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
		for _, photo := range photos {
			presignedUrl := utils.GeneratePresignedUrl(os.Getenv("FILE_UPLOAD_PATH")+photo.BaseUrl, os.Getenv("SIGNED_URL_SECRET"), 1*time.Hour) // Adjust the parameters as needed
			photoMap := utils.StructToMap(photo)
			photoMap["PresignedUrl"] = scheme + "://" + c.Request.Host + "/api/public/files/" + presignedUrl
			photoWithPresignedUrl = append(photoWithPresignedUrl, photoMap)
		}

		c.JSON(http.StatusOK, types.SuccessSearchResponse{Data: photoWithPresignedUrl, Page: searchParams.Page, PageSize: searchParams.PerPage, Total: count})
	} else {
		c.JSON(http.StatusOK, types.SuccessSearchResponse{Data: photos, Page: searchParams.Page, PageSize: searchParams.PerPage, Total: count})
	}
	return
}

func (ctrl PhotoHandlerV1) Delete(c *gin.Context) {
	//delete photo
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	userID := introspection.Sub
	photo := models.Photo{}
	user := models.User{}
	filter := map[string]any{
		"id": c.Param("photo_id"),
	}
	searchParams := db.SearchParams{
		Filters: filter,
	}
	db.QueryOne(&searchParams, nil, &photo)
	if photo.ID == 0 {
		c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "Photo not found"})
		return
	}

	db.QueryOne(&db.SearchParams{Filters: map[string]any{"id": userID}, Populate: []string{"Albums"}}, nil, &user)
	//map user.Albums to albumIDs
	photoInAlbum := false
	for _, album := range user.Albums {
		if album.ID == photo.AlbumID {
			photoInAlbum = true
			break
		}
	}
	if !photoInAlbum {
		c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "Ypu don't have permission to delete this photo"})
		return
	}
	result := db.DB.Delete(&photo)
	//delete from storage
	if os.Getenv("UPLOAD_DRIVER") == "local" {
		os.Remove(os.Getenv("FILE_UPLOAD_PATH") + photo.BaseUrl)
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, types.Error{Code: http.StatusInternalServerError, Message: "Error deleting photo"})
		return
	}
	c.JSON(http.StatusOK, types.SuccessResponse{Data: photo})
	return
}

func (ctrl PhotoHandlerV1) Update(c *gin.Context) {
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	userID := introspection.Sub
	user := models.User{}
	updatePhotoValidator := validatorsV1.UpdatePhotoValidator{}
	if err := c.ShouldBindJSON(&updatePhotoValidator); err != nil {
		c.JSON(http.StatusNotAcceptable, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusNotAcceptable,
			Message: err.Error(),
			Details: err,
		}})
		c.Abort()
		return
	}

	photoID := c.Param("photo_id")
	photo := models.Photo{}
	var searchParams = db.SearchParams{
		Filters: map[string]any{
			"id": photoID,
		},
	}
	db.QueryOne(&searchParams, nil, &photo)
	if photo.ID == 0 {
		c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "Photo not found"})
		return
	}
	db.QueryOne(&db.SearchParams{Filters: map[string]any{"id": userID}, Populate: []string{"Albums"}}, nil, &user)
	photoInAlbum := false
	for _, album := range user.Albums {
		if album.ID == photo.AlbumID {
			photoInAlbum = true
			break
		}
	}
	if !photoInAlbum {
		c.JSON(http.StatusForbidden, types.Error{Code: http.StatusForbidden, Message: "Ypu don't have permission to delete this photo"})
		return
	}
	photo.BaseName = updatePhotoValidator.Name
	result := db.DB.Save(&photo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, types.Error{Code: http.StatusInternalServerError, Message: "Error updating photo"})
		return
	}
	c.JSON(http.StatusOK, types.SuccessResponse{Data: photo})
}
