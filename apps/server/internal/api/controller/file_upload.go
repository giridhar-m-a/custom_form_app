package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func FileUploadController(rg *gin.RouterGroup) {
	api := rg.Group("/files")
	api.POST("/file-upload", handler.FileUploadHandler)
}
