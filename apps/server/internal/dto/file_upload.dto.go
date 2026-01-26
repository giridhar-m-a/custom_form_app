package dto

import "mime/multipart"

type FileUploadPayload struct {
	File     multipart.File
	FileInfo *multipart.FileHeader
}
