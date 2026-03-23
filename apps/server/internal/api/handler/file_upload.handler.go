package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

// FileUploadHandler godoc
// @Summary Upload a file to MinIO
// @Description Uploads a file to the specified MinIO bucket and returns file info
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param path formData string true "File path to upload shouldbe formid/invitationid"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /files/file-upload [post]
// @Security BearerAuth
// @type http
// @scheme bearer
func FileUploadHandler(c *gin.Context) {
	// Parse the multipart form, with a max memory of 10MB
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Retrieve the file from form data
	file, header, err := c.Request.FormFile("file")
	path := c.PostForm("path")
	if err != nil {
		c.JSON(400, gin.H{"message": "Failed to retrieve file", "status": http.StatusBadRequest})
		return
	}
	defer file.Close()

	// Upload the file to MinIO
	bucketName := utils.GetEnv("MINIO_BUCKET_NAME", "custom-form-app") // Replace with your bucket name
	objectName := header.Filename
	contentType := header.Header.Get("Content-Type")

	objectPath := fmt.Sprintf("%s/%s", path, objectName)

	fileInfo, err := services.MinioUploadFile(bucketName, objectPath, file, header.Size, contentType)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	signedUrl, err := services.GetMinioSignedURL(bucketName, objectName, time.Hour*24, fileInfo.VersionID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(200, gin.H{
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"file_info":  fileInfo,
			"objectName": objectName,
			"signed":     signedUrl,
		},
		"status": http.StatusOK,
	})
}
