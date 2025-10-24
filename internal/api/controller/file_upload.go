package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

// FileUploadHandler godoc
// @Summary Upload a file to MinIO
// @Description Uploads a file to the specified MinIO bucket and returns file info
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /file-upload [post]
// @Security BearerAuth
// @type http
// The route expects multipart/form-data, produces JSON responses, and is secured with Bearer authentication.
func FileUploadController(rg *gin.RouterGroup) {
	rg.POST("/file-upload", handler.FileUploadHandler)
}