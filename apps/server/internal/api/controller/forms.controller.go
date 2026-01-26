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
		forms.GET("/", formHandler.GetForms)
		forms.GET("/:formID", formHandler.GetSingleForm)
		forms.PATCH("/:formID", formHandler.UpdateForm)
		forms.DELETE("/:formID", formHandler.DeleteForm)
		forms.GET("/fields/:formID", formHandler.GetFormFields)
		forms.PATCH("/fields", formHandler.UpdateFormFields)
	}
}
