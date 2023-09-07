package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/auditt98/onthego/db"
	"github.com/auditt98/onthego/models"
	"github.com/auditt98/onthego/types"
	"github.com/auditt98/onthego/utils"
	"github.com/gin-gonic/gin"
)

type PhotoHandlerV1 struct{}

func (ctrl PhotoHandlerV1) Search(c *gin.Context) {
	//get photos
	introspection, _ := c.Get("introspectionResult")
	userID := introspection.(*types.IntrospectionResult).Sub
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
			photoMap["PresignedUrl"] = scheme + "://" + c.Request.Host + "/api/v1/files/" + presignedUrl
			photoWithPresignedUrl = append(photoWithPresignedUrl, photoMap)
		}

		c.JSON(http.StatusOK, types.SuccessSearchResponse{Data: photoWithPresignedUrl, Page: searchParams.Page, PageSize: searchParams.PerPage, Total: count})
	} else {
		c.JSON(http.StatusOK, types.SuccessSearchResponse{Data: photos, Page: searchParams.Page, PageSize: searchParams.PerPage, Total: count})
	}
	return
}
