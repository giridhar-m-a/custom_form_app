package api

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller/auth"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

// RegisterRoutes registers the API v1 routes on the provided Gin engine.
// It mounts routes under the "/api/v1" prefix (health and Google auth), then applies the authentication middleware and registers protected routes such as file upload and forms.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	controller.RegisterHealth(api)
	auth.GoogleAuth(api)
	api.Use(middleware.AuthMiddleware())
	controller.FileUploadController(api)
	controller.RegisterFormsController(api)
}