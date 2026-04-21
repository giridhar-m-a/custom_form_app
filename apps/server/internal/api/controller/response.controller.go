package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

func RegisterResponseRoutes(router *gin.RouterGroup) {
	responseHandler := handler.NewResponseHandler()

	rg := router.Group("/response")
	rg.POST("/", responseHandler.CreateSubmission)
	protected := rg.Use(middleware.AuthMiddleware())
	protected.GET("/:formId", responseHandler.GetSubmissions)
	protected.GET("/submission/:submissionId", responseHandler.GetSingleSubmission)
}
