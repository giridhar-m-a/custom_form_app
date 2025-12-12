package dto

import "time"

// API response DTO
type User struct {
	UserID           string    `json:"id"`
	UserEmail        string    `json:"email"`
	UserFullName     string    `json:"fullName"`
	UserProfilePic   *string   `json:"profilePic"`
	UserProfilePicID *string   `json:"profilePicId"`
	UserCreatedAt    time.Time `json:"createdAt"`
	UserUpdatedAt    time.Time `json:"updatedAt"`
}
