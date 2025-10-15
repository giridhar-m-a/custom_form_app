package controller

import "github.com/gin-gonic/gin"

func RegisterFormsController(rg *gin.RouterGroup) {
	forms := rg.Group("/forms")
	{
		forms.GET("/:id", getFormByID)
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
// @scheme bearer
func getFormByID(c *gin.Context) {
	formID := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Form retrieved successfully",
		"formID":  formID,
		})
}
