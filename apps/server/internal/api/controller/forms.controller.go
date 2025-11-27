package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

// RegisterFormsController sets up form routes
func RegisterFormsController(rg *gin.RouterGroup) {
	forms := rg.Group("/forms")

	// Create a handler instance
	formHandler := handler.NewFormsHandler()

	{
		forms.POST("/", formHandler.CreateForm)
		forms.POST("/fields", formHandler.CreateFormFields)
	}
}