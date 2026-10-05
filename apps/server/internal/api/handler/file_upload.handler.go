package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
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
// @Success 200 {object} dto.ApiResponse[dto.FileUploadResponse] "file uploaded successful"
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

	signedUrl, err := services.GetMinioSignedURL(bucketName, objectPath, time.Hour*24, fileInfo.VersionID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(200, dto.ApiResponse[dto.FileUploadResponse]{
		Data: dto.FileUploadResponse{
			FileInfo: dto.FileResponse{
				FilePath: objectPath,
				FileName: objectName,
				FileType: header.Header.Get("Content-Type"),
				FileSize: header.Size,
			},
			Signed: signedUrl.String(),
		},
		Status:  200,
		Message: "file uploaded successful",
	})
}



// @Summary Get signed url of file
// @Description Get signed url of file
// @Tags Files
// @Accept       json
// @Produce      json
// @Param        form  body      dto.GetSignedUrlPayload  true  "Form data"
// @Success      201   {object}  dto.ApiResponse[dto.GetSignedUrlResponse]  "Form created successfully"
// @Failure      400   {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401   {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500   {object}  object{status=string,message=string}  "Internal server error"
// @Router /files/get-signed-url [post]
// @Security BearerAuth
// @type http
// @scheme bearer
func GetSignedUrl(c *gin.Context) {
	var payload dto.GetSignedUrlPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.HandleError(c,err)
		return
	}

	// Upload the file to MinIO
	bucketName := utils.GetEnv("MINIO_BUCKET_NAME", "custom-form-app") // Replace with your bucket name
	

	signedUrl, err := services.GetMinioSignedURL(bucketName, payload.FilePath, time.Hour*24, "")
	if err != nil {
		utils.HandleError(c,err)
		return
	}

	c.JSON(200, dto.ApiResponse[dto.GetSignedUrlResponse]{
		Data: dto.GetSignedUrlResponse{
			SignedUrl: signedUrl.String(),
		},
		Status:  200,
		Message: "file uploaded successful",
	})
}