package handlers

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	"github.com/auditt98/onthego/utils"
	validatorsV1 "github.com/auditt98/onthego/validators/v1"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

type AlbumHandlerV1 struct{}

// @Summary Create an album
// @Description Create a new album with the provided data.
// @Tags Albums
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT access token"
// @Param albumValidator body  validatorsV1.NewAlbumValidator{} true "Album data to be created"
// @Success 200 {object} types.SuccessResponse{Data=models.Album} "Successful response with the created album"
// @Failure 406 {object} types.ErrorResponse{Error=types.Error} "Validation error or invalid data"
// @Failure 500 {object} types.ErrorResponse{Error=types.Error} "Internal server error"
// @Router /api/v1/albums [post]
func (ctrl AlbumHandlerV1) CreateAlbum(c *gin.Context) {
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
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
		ID: introspection.Sub,
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
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
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
		if user.ID == introspection.Sub {
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
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	userID := introspection.Sub
	albums := []models.Album{}
	currentUserFilter := map[string]any{
		"users": map[string]any{
			"id": userID,
		},
	}
	var count int64
	searchParams := db.SearchParams{}
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		db.DB.Where(currentUserFilter).Find(&albums)
	}
	db.Query(&searchParams, currentUserFilter, &albums, &count)
	isPopulatePhotos := slices.ContainsFunc(searchParams.Populate, func(s string) bool {
		return strings.Contains(s, "Photos")
	})

	var scheme string
	// c.Request.Header.Get("X-Forwarded-Proto")
	//check if c.Request.Header.Get("X-Forwarded-Proto") exists, if it does, use that, else check if c.Request.TLS is true, if it is, use https, else use http
	if c.Request.Header.Get("X-Forwarded-Proto") != "" {
		scheme = c.Request.Header.Get("X-Forwarded-Proto")
	} else if c.Request.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	if isPopulatePhotos && os.Getenv("UPLOAD_DRIVER") == "local" {
		for _, album := range albums {
			photoWithPresignedUrl := []*models.Photo{}
			for _, photo := range album.Photos {
				presignedUrl := utils.GeneratePresignedUrl(os.Getenv("FILE_UPLOAD_PATH")+photo.BaseUrl, os.Getenv("SIGNED_URL_SECRET"), 1*time.Hour)
				photo.PresignedUrl = scheme + "://" + c.Request.Host + "/api/public/files/" + presignedUrl
				photoWithPresignedUrl = append(photoWithPresignedUrl, photo)
			}
			album.Photos = photoWithPresignedUrl
		}
	}
	c.JSON(http.StatusOK, types.SuccessSearchResponse{Data: albums, Page: searchParams.Page, PageSize: searchParams.PerPage, Total: count})
	return
}

func (ctrl AlbumHandlerV1) RemoveUser(c *gin.Context) {
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))

	albumID := c.Param("album_id")
	userID := c.Param("user_id")

	var resultAlbum models.Album
	db.DB.Preload("Users").First(&resultAlbum, albumID)
	var isCurrentUserInAlbum bool
	var isUserAlreadyInAlbum bool

	for _, user := range resultAlbum.Users {
		if user.ID == introspection.Sub {
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
		if user.ID == introspection.Sub {
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
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	userID := introspection.Sub
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

	form, err := c.MultipartForm()
	if (err != nil) || (form == nil) {
		c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	files := form.File["files[]"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "No files found"})
		return
	}
	paths, err := utils.FileUpload(files, albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	photos := []*models.Photo{}

	for _, file := range paths {
		//create Photo, baseName would be file.split('/')[-1] then url decode
		baseName := path.Base(file)
		baseName, _ = url.QueryUnescape(baseName)

		photo := models.Photo{
			BaseName: baseName,
			BaseUrl:  file,
			AlbumID:  albums[0].ID,
			UserID:   userID,
		}
		photos = append(photos, &photo)
	}
	result := db.DB.Create(&photos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, types.Error{Code: http.StatusInternalServerError, Message: result.Error.Error()})
		c.Abort()
		return
	}
	var scheme string
	// c.Request.Header.Get("X-Forwarded-Proto")
	//check if c.Request.Header.Get("X-Forwarded-Proto") exists, if it does, use that, else check if c.Request.TLS is true, if it is, use https, else use http
	if c.Request.Header.Get("X-Forwarded-Proto") != "" {
		scheme = c.Request.Header.Get("X-Forwarded-Proto")
	} else if c.Request.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	if os.Getenv("UPLOAD_DRIVER") == "local" {
		photoWithPresignedUrl := []*models.Photo{}
		for _, photo := range photos {
			presignedUrl := utils.GeneratePresignedUrl(os.Getenv("FILE_UPLOAD_PATH")+photo.BaseUrl, os.Getenv("SIGNED_URL_SECRET"), 1*time.Hour) // Adjust the parameters as needed
			photo.PresignedUrl = scheme + "://" + c.Request.Host + "/api/public/files/" + presignedUrl
			photoWithPresignedUrl = append(photoWithPresignedUrl, photo)
		}
		c.JSON(http.StatusOK, types.SuccessResponse{Data: photoWithPresignedUrl})
	} else {
		c.JSON(http.StatusOK, types.SuccessResponse{Data: photos})
	}
	return
}

func (ctrl AlbumHandlerV1) UpdateAlbum(c *gin.Context) {
	introspection := utils.GetIntrospection(c.Get("introspectionResult"))
	updateAlbumValidator := validatorsV1.UpdateAlbumValidator{}
	if err := c.ShouldBindJSON(&updateAlbumValidator); err != nil {
		c.JSON(http.StatusNotAcceptable, types.ErrorResponse{Error: types.Error{
			Code:    http.StatusNotAcceptable,
			Message: err.Error(),
			Details: err,
		}})
		c.Abort()
		return
	}

	albumID := c.Param("album_id")
	album := models.Album{}
	var searchParams = db.SearchParams{
		Filters: map[string]any{
			"id": albumID,
			"users": map[string]any{
				"id": introspection.Sub,
			},
		},
	}
	db.QueryOne(&searchParams, nil, &album)
	if album.ID == 0 {
		c.JSON(http.StatusNotFound, types.Error{Code: http.StatusNotFound, Message: "Album not found"})
		return
	}
	album.Name = updateAlbumValidator.Name
	db.DB.Save(&album)
	c.JSON(http.StatusOK, types.SuccessResponse{Data: album})
	return
}
