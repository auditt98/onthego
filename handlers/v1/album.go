package handlers

import (
	"net/http"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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

func (ctrl AlbumHandlerV1) AddUserToAlbum(c *gin.Context) {
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
	searchParams := db.SearchParams{}
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		db.DB.Where(currentUserFilter).Find(&albums)
	}
	var count int64
	db.DB.Where(searchParams.Filters).Where(currentUserFilter).Preload(clause.Associations).Scopes(db.Paginate(&searchParams.Page, &searchParams.PerPage)).Find(&albums).Count(&count)
	c.JSON(http.StatusOK, types.SuccessResponse{Data: albums, Page: searchParams.Page, PageSize: searchParams.PerPage, Total: count})
	return
}

func (ctrl AlbumHandlerV1) RemoveUserFromAlbum(c *gin.Context) {
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
