package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

// RegisterFormsController registers routes under the /forms path on the provided router group.
// It attaches GET /forms/:id to retrieve a form by ID and POST /forms/ to create a new form.
func RegisterFormsController(rg *gin.RouterGroup) {
	forms := rg.Group("/forms")

	// Create a handler instance
	formHandler := handler.NewFormsHandler()

	{
		forms.GET("/:id", getFormByID)
		forms.POST("/", formHandler.CreateForm) // ✅ Use the handler method here
	}
}

// GetFormByID retrieves a form by its ID
// @Summary      Get a form by ID
// @Description  Retrieves a form with the specified ID
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Form ID"
// @Success      200  {object}  object{message=string,formID=string}  "Form retrieved successfully"
// @Failure      400  {object}  object{status=int,message=string}     "Bad request"
// @Failure      401  {object}  object{status=int,message=string}     "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}     "Internal server error"
// @Schemes      https
// @Router       /forms/{id} [get]
// @Security BearerAuth
// @type http
// getFormByID writes a 200 JSON response containing a success message and the requested form ID from the path parameter `id`.
//
// The response JSON contains the keys "message" and "formID". The form ID is taken directly from the request path parameter "id".
func getFormByID(c *gin.Context) {
	formID := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Form retrieved successfully",
		"formID":  formID,
	})
}