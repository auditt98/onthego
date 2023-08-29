package handlers

import (
	"fmt"
	"net/http"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
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
		Id: introspection.(*types.IntrospectionResult).Sub,
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
	fmt.Println("albumID", albumID)
	var resultAlbums []models.Album
	db.DB.Table("albums").InnerJoins("users").Find(&resultAlbums, map[string]interface{}{
		"id": albumID,
	})

	// db.DB.Table("albums").Preload("Users", func(d *gorm.DB) *gorm.DB {
	// 	return d.Where("id = ?", introspection.(*types.IntrospectionResult).Sub)
	// }).Find(&resultAlbums, map[string]interface{}{
	// 	"id": albumID,
	// })

	db.DB.Table("albums").Preload("Users").Where("id = ? and Users.id = ?", albumID, introspection.(*types.IntrospectionResult).Sub).Find(&resultAlbums)

	// db.DB.Preload("Users", {
	// 	"id": introspection.(*types.IntrospectionResult).Sub,
	// }).Find(&resultAlbums)

	// params := db.QueryParams{
	// 	Where: db.WhereParams{
	// 		And: []db.WhereParams{
	// 			{
	// 				Attr: map[string]db.AttributeParams{
	// 					"id": {
	// 						Eq: albumID,
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Attr: map[string]db.AttributeParams{
	// 					"Users.id": {
	// 						Contains: introspection.(*types.IntrospectionResult).Sub,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	OrderBy:  []string{"-id"},
	// 	Limit:    1,
	// 	Offset:   0,
	// 	Populate: []string{"Users"},
	// }

	// albums := []models.Album{}
	// db.Query("albums", params, &albums)
	// fmt.Println("albums", albums)

	// if result.Error.Error() != "" {
	// 	c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: types.Error{
	// 		Code:    http.StatusInternalServerError,
	// 		Message: "Error adding user to album",
	// 		Details: result.Error.Error(),
	// 	}})
	// 	c.Abort()
	// 	return
	// }
	c.JSON(http.StatusOK, types.SuccessResponse{Data: resultAlbums})

	// var album models.Album

	// album = models.Album{
	// 	ID: addUserToAlbumValidator.AlbumID,
	// }

	// for _, userId := range addUserToAlbumValidator.UserIds {
	// 	album.Users = append(album.Users, &models.User{
	// 		Id: userId,
	// 	})
	// }

	// result := db.DB.Save(&album)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: types.Error{
	// 		Code:    http.StatusInternalServerError,
	// 		Message: "Error adding user to album",
	// 		Details: result.Error,
	// 	}})
	// 	c.Abort()
	// 	return
	// }
	// c.JSON(http.StatusOK, types.SuccessResponse{Data: album})

	return
}
