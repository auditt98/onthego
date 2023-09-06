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

type AlbumHandlerV1 struct{}

func (ctrl AlbumHandlerV1) CreateAlbum(c *gin.Context) {
	introspection, _ := c.Get("introspectionResult")
	albumValidator := validatorsV1.NewAlbumValidator{}
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
	user := models.User{
		ID: introspection.(*types.IntrospectionResult).Sub,
	}
	newAlbum.Name = albumValidator.Name
	newAlbum.Users = append(newAlbum.Users, &user)
	result := db.DB.Create(&newAlbum)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error creating album",
			Details: result.Error,
		}})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, types.SuccessResponse{Data: newAlbum})
	return
}

func (ctrl AlbumHandlerV1) AddUser(c *gin.Context) {
	introspection, _ := c.Get("introspectionResult")
	addUserToAlbumValidator := validatorsV1.AddUserToAlbumValidator{}
	if err := c.ShouldBindJSON(&addUserToAlbumValidator); err != nil {
		c.JSON(http.StatusNotAcceptable, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusNotAcceptable,
			Message: err.Error(),
			Details: err,
		}})
		c.Abort()
		return
	}
	albumID := c.Param("album_id")
	var resultAlbum models.Album
	db.DB.Preload("Users").First(&resultAlbum, albumID)
	var isCurrentUserInAlbum bool
	var isUserAlreadyInAlbum bool

	for _, user := range resultAlbum.Users {
		if user.ID == introspection.(*types.IntrospectionResult).Sub {
			isCurrentUserInAlbum = true
		}
		if user.ID == addUserToAlbumValidator.UserId {
			isUserAlreadyInAlbum = true
		}
	}

	if !isCurrentUserInAlbum {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusForbidden,
			Message: "Current user does not belong to the album",
		}})
		c.Abort()
		return
	}

	if isUserAlreadyInAlbum {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusForbidden,
			Message: "User already belongs to the album",
		}})
		c.Abort()
		return
	}

	resultAlbum.Users = append(resultAlbum.Users, &models.User{ID: addUserToAlbumValidator.UserId})
	db.DB.Save(&resultAlbum)
	err := db.DB.Save(&resultAlbum).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, types.SuccessResponse{Data: resultAlbum})
	return
}

func (ctrl AlbumHandlerV1) Search(c *gin.Context) {
	introspection, _ := c.Get("introspectionResult")
	userID := introspection.(*types.IntrospectionResult).Sub
	albums := []models.Album{}
	currentUserFilter := map[string]any{
		"users": map[string]any{
			"id": userID,
		},
	}
	var count int64
	searchParams2 := db.SearchParams{}
	if err := c.ShouldBindJSON(&searchParams2); err != nil {
		db.DB.Where(currentUserFilter).Find(&albums)
	}
	db.Query(&searchParams2, currentUserFilter, &albums, &count)
	c.JSON(http.StatusOK, types.SuccessSearchResponse{Data: albums, Page: searchParams2.Page, PageSize: searchParams2.PerPage, Total: count})
	return
}

func (ctrl AlbumHandlerV1) RemoveUser(c *gin.Context) {
	introspection, _ := c.Get("introspectionResult")

	albumID := c.Param("album_id")
	userID := c.Param("user_id")

	var resultAlbum models.Album
	db.DB.Preload("Users").First(&resultAlbum, albumID)
	var isCurrentUserInAlbum bool
	var isUserAlreadyInAlbum bool

	for _, user := range resultAlbum.Users {
		if user.ID == introspection.(*types.IntrospectionResult).Sub {
			isCurrentUserInAlbum = true
		}
		if user.ID == userID {
			isUserAlreadyInAlbum = true
		}
	}

	if !isCurrentUserInAlbum {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusForbidden,
			Message: "Current user does not belong to the album",
		}})
		c.Abort()
		return
	}

	if !isUserAlreadyInAlbum {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusForbidden,
			Message: "User does not belong to the album",
		}})
		c.Abort()
		return
	}

	var updatedUsers []*models.User
	userIDToRemove := userID
	for _, user := range resultAlbum.Users {
		if user.ID != userIDToRemove {
			updatedUsers = append(updatedUsers, user)
		}
	}

	resultAlbum.Users = updatedUsers
	if err := db.DB.Save(&resultAlbum).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}})
		c.Abort()
		return
	}
	//if current user is not in the album anymore,
	var currentUserInUpdatedAlbum bool
	for _, user := range resultAlbum.Users {
		if user.ID == introspection.(*types.IntrospectionResult).Sub {
			currentUserInUpdatedAlbum = true
			break
		}
	}

	if !currentUserInUpdatedAlbum {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusForbidden,
			Message: "Current user is no longer in the album",
		}})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, types.SuccessResponse{Data: resultAlbum})
}

func (ctrl AlbumHandlerV1) AddPhotos(c *gin.Context) {
	introspection, _ := c.Get("introspectionResult")
	userID := introspection.(*types.IntrospectionResult).Sub
	albumID := c.Param("album_id")
	albums := []models.Album{}
	var searchParams = db.SearchParams{
		Filters: map[string]any{
			"id": albumID,
			"users": map[string]any{
				"id": userID,
			},
		},
		Page:    1,
		PerPage: 1,
	}
	var count int64
	db.Query(&searchParams, nil, &albums, &count)
	if count == 0 || len(albums) == 0 {
		c.JSON(http.StatusNotFound, types.Error{Code: http.StatusNotFound, Message: "Album not found"})
	}

	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "No files found"})
		return
	}
	utils.FileUpload(files, albumID)
	for _, file := range files {
		c.JSON(http.StatusOK, types.SuccessResponse{Data: file})
	}

	return
}

func (ctrl AlbumHandlerV1) SearchPhotos(c *gin.Context) {
	return
}
