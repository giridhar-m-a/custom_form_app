package dto

import "mime/multipart"

type FileUploadPayload struct {
	File     multipart.File
	FileInfo *multipart.FileHeader
}

type FileResponse struct {
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
}

type SignedUrlResponse struct {
	Scheme      string `json:"scheme"`
	Opaque      string `json:"opaque"`
	User        string `json:"user"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Fragment    string `json:"fragment"`
	RawQuery    string `json:"rawQuery"`
	RawPath     string `json:"rawPath"`
	RawFragment string `json:"rawFragment"`
	ForceQuery  bool   `json:"forceQuery"`
	OmitHost    bool   `json:"omitHost"`
}

type FileUploadResponse struct {
	FileInfo FileResponse `json:"fileInfo"`
	Signed   string       `json:"signedUrl"`
}

type GetSignedUrlPayload struct {
	FilePath string `json:"filePath"`
}

type GetSignedUrlResponse struct {
	SignedUrl string `json:"signedUrl"`
}