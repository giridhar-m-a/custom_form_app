package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

// RegisterFormsController sets up form routes
func RegisterFormsController(rg *gin.RouterGroup) {
	responseForms := rg.Group("/forms")
	forms := rg.Group("/forms")

	// Create a handler instance
	formHandler := handler.NewFormsHandler()

	response := responseForms.Use(middleware.ResponseMiddleware())
	{
		response.GET("/response", formHandler.GetFormForResponse)
		response.GET("/fields/response", formHandler.GetFormFieldsForResponse)
	}
	protected := forms.Use(middleware.AuthMiddleware())

	{
		protected.POST("/", formHandler.CreateForm)
		protected.POST("/fields", formHandler.CreateFormFields)
		protected.GET("/", formHandler.GetForms)
		protected.GET("/:formID", formHandler.GetSingleForm)
		protected.PATCH("/:formID", formHandler.UpdateForm)
		protected.DELETE("/:formID", formHandler.DeleteForm)
		protected.GET("/fields/:formID", formHandler.GetFormFields)
		protected.PATCH("/fields", formHandler.UpdateFormFields)
	}
}
