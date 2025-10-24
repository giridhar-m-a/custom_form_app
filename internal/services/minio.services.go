package services

import (
	"context"
	"errors"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient is the shared MinIO client
var MinioClient *minio.Client

// InitMinio initializes the package-level MinioClient using MinIO configuration from environment variables.
// It reads MINIO_SERVER, MINIO_USER, MINIO_PASSWORD and MINIO_USE_SSL, attempts to create a MinIO client and assigns it to MinioClient.
// If client creation fails the error is logged and MinioClient will be nil.
func InitMinio() {
	endpoint := utils.GetEnv("MINIO_SERVER", "")
	accessKey := utils.GetEnv("MINIO_USER", "")
	secretKey := utils.GetEnv("MINIO_PASSWORD", "")
	useSSL := utils.GetEnvAsBool("MINIO_USE_SSL", false)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to initialize MinIO client: %v", err)
	}

	MinioClient = client
	log.Println("MinIO client initialized successfully")
}

// UploadFile uploads a file to the specified bucket
func MinioUploadFile(bucketName, objectName string, reader io.Reader, size int64, contentType string) (minio.UploadInfo, error) {
	if MinioClient == nil {
		log.Println("MinIO client not initialized")
		return minio.UploadInfo{}, errors.New("MinIO client not initialized")
	}

	fileRes, err := MinioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		log.Println("Failed to upload file to MinIO:", err)
		return minio.UploadInfo{}, err
	}

	return fileRes, nil
}

// GetPresignedURL generates a signed URL valid for `expiry` duration
func GetMinioSignedURL(bucketName, objectName string, expiry time.Duration, versionID string) (*url.URL, error) {
	if MinioClient == nil {
		log.Fatalln("MinIO client not initialized")
	}

	reqParams := make(url.Values) // additional query params if needed
	if versionID != "" {
		reqParams.Set("versionId", versionID)
	}
	presignedURL, err := MinioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}