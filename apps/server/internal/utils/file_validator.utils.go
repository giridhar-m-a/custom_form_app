package utils

import (
	"mime/multipart"
	"net/http"
)

func IsValidImage(file multipart.File) (bool, string, error) {
	defer file.Seek(0, 0) // reset pointer

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false, "", err
	}

	contentType := http.DetectContentType(buffer)

	switch contentType {
	case "image/jpeg", "image/png", "image/webp", "image/gif":
		return true, contentType, nil
	default:
		return false, contentType, nil
	}
}
