package handlers

import "github.com/gin-gonic/gin"

type FileHandlerV1 struct{}

func (ctrl FileHandlerV1) GetFile(c *gin.Context) {
	filePath := c.Param("file_path")
	filePath = filePath[1:]
	c.File(filePath)
	return
}
