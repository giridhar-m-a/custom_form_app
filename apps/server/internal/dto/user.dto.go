package dto

import "time"

// API response DTO
type User struct {
	UserID         string    `json:"id"`
	UserEmail      string    `json:"email"`
	UserFullName   string    `json:"fullName"`
	UserCreatedAt  time.Time `json:"createdAt"`
	UserUpdatedAt  time.Time `json:"updatedAt"`
	UserProfilePic string    `json:"profilePic"`
	IsTemp         bool      `json:"isTemp"`
}

type UserProfilePicResponse struct {
	UserID   string `json:"userId"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
}

type UserUpdateDTO struct {
	UserFullName       string `json:"userFullName"`
	UserPassword       string `json:"userPassword"`
	UserVerifyPassword string `json:"userVerifyPassword"`
	UserGoogleId       string `json:"userGoogleId"`
}

type UpdateUserDetailsDTO struct {
	UserFullName string `json:"userFullName"`
}

type UpdateUserPasswordDTO struct {
	OldPassword        string `json:"oldPassword" binding:"required"`
	UserPassword       string `json:"userPassword" binding:"required"`
	UserVerifyPassword string `json:"userVerifyPassword" binding:"required"`
}
