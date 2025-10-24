package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
)

// and "status".
func FileUploadHandler(c *gin.Context) {
	// Parse the multipart form, with a max memory of 10MB
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Retrieve the file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message": "Failed to retrieve file", "status": http.StatusBadRequest})
		return
	}
	defer file.Close()

	// Upload the file to MinIO
	bucketName := "user-image" // Replace with your bucket name
	objectName := header.Filename
	contentType := header.Header.Get("Content-Type")

	fileInfo, err := services.MinioUploadFile(bucketName, objectName, file, header.Size, contentType)
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
		"message":    "File uploaded successfully",
		"file_info":  fileInfo,
		"bucket":     bucketName,
		"objectName": objectName,
		"signedUrl":  signedUrl.String(),
		"signed":     signedUrl,
		"status":     http.StatusOK,
	})
}