package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/auditt98/onthego/types"
	"github.com/auditt98/onthego/utils"
	"github.com/gin-gonic/gin"
)

func PresignedUrlValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// // Get the presigned URL from the request
		fmt.Println("Presigned URL Validator")
		filePath := c.Param("file_path")
		expires := c.DefaultQuery("expires", "")
		signature := c.DefaultQuery("signature", "")
		if filePath == "" || expires == "" || signature == "" {
			c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "Invalid presigned URL"})
			c.Abort()
			return
		}
		// Validate the presigned URL
		expiresTime, err := strconv.ParseInt(expires, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "Invalid presigned URL"})
			c.Abort()
			return
		}
		currentTime := time.Now().Unix()
		if expiresTime <= currentTime {
			c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "Invalid presigned URL: expires timestamp has expired"})
			c.Abort()
			return
		}
		//remove first / from filePath
		filePath = filePath[1:]
		generatedSignature := utils.CalculateSignature(filePath, os.Getenv("SIGNED_URL_SECRET"), expiresTime)
		if generatedSignature != signature {
			c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "Invalid presigned URL: signature does not match"})
			c.Abort()
			return
		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, types.Error{Code: http.StatusBadRequest, Message: "File not found"})
			c.Abort()
			return
		}
		c.Next()
	}

}
