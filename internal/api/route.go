package api

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller/auth"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	controller.RegisterHealth(api)
	auth.GoogleAuth(api)
	controller.FileUploadController(api)
}
