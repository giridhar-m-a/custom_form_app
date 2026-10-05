package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

func FileUploadController(rg *gin.RouterGroup) {
	api := rg.Group("/files")
	protected:= rg.Group("/files")
	api.Use(middleware.ResponseMiddleware())
	api.POST("/file-upload", handler.FileUploadHandler)
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/get-signed-url", handler.GetSignedUrl)
}
