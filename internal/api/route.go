package api

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/controller"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	controller.RegisterHealth(api)
}
