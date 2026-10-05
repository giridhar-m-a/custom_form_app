package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient is the shared MinIO client
var MinioClient *minio.Client

// InitMinio initializes the MinIO client
func InitMinio() {
	// endpoint := utils.GetEnv("MINIO_SERVER", "")
	accessKey := utils.GetEnv("MINIO_USER", "")
	secretKey := utils.GetEnv("MINIO_PASSWORD", "")
	useSSL := utils.GetEnvAsBool("MINIO_USE_SSL", false)
	domain:=utils.GetEnv("MINIO_DOMAIN","")
	client, err := minio.New(domain, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to initialize Storage: %v", err)
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

func MinioDeleteFile(bucketName, objectName string) error {
	if MinioClient == nil {
		log.Println("MinIO client not initialized")
		return errors.New("MinIO client not initialized")
	}

	err := MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Failed to delete file from MinIO:", err)
		return err
	}

	return nil
}

// GetPresignedURL generates a signed URL valid for `expiry` duration
func GetMinioSignedURL(bucketName, objectName string, expiry time.Duration, versionID string) (*url.URL, error) {
	if MinioClient == nil {
		log.Printf("MinIO client not initialized")
		return nil, errors.New("MinIO client not initialized")
	}
	minioClientDomain := utils.GetEnv("MINIO_DOMAIN", "minio.custom-form-app.home")
	reqParams := make(url.Values) // additional query params if needed
	if versionID != "" {
		reqParams.Set("versionId", versionID)
		reqParams.Set("Host", minioClientDomain)
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

func DeleteFolderBulk(bucketName, folderPrefix string, ctx context.Context) error {
	if MinioClient == nil {
		log.Printf("MinIO client not initialized")
		return errors.New("MinIO client not initialized")
	}
	if !strings.HasSuffix(folderPrefix, "/") {
		folderPrefix += "/"
	}

	objectsCh := MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    folderPrefix,
		Recursive: true,
	})

	removeObjectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(removeObjectsCh)
		for obj := range objectsCh {
			if obj.Err != nil {
				log.Printf("[MinIO] Error listing object: %v", obj.Err)
				continue
			}
			log.Printf("[MinIO] Queuing for deletion: %s", obj.Key) // 👈 log key before sending
			removeObjectsCh <- obj
		}
	}()

	for removeErr := range MinioClient.RemoveObjects(ctx, bucketName, removeObjectsCh, minio.RemoveObjectsOptions{}) {
		if removeErr.Err != nil {
			return fmt.Errorf("bulk delete error on object %s: %w", removeErr.ObjectName, removeErr.Err) // 👈 ObjectName available here
		}
	}

	log.Printf("[MinIO] Successfully deleted all objects under prefix: %s", folderPrefix)
	return nil
}