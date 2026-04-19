package api

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.Use(middleware.UnknownFieldsMiddleware())
	controller.RegisterHealth(api)
	controller.Auth(api)
	controller.RegisterResponseRoutes(api)
	controller.RegisterInvitationsRoutes(api)
	controller.RegisterFormsController(api)
	controller.FileUploadController(api)
	api.Use(middleware.AuthMiddleware())

	controller.Users(api)
}
