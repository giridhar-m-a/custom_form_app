package api

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	controller.RegisterHealth(api)
	controller.Auth(api)
	api.Use(middleware.AuthMiddleware())
	controller.FileUploadController(api)
	controller.RegisterFormsController(api)
	controller.Users(api)
}
