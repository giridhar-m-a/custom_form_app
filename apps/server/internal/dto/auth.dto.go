package dto

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

type GoogleAuthRequest struct {
	Code string `json:"code" binding:"required" form:"code" message:"code is required"`
}

type EmailPasswordAuthRequest struct {
	Email    string `json:"email" binding:"required,email" message:"valid email is required"`
	Password string `json:"password" binding:"required,min=6" message:"password with minimum 6 characters is required"`
}

type EmailPasswordRegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2" message:"full name with minimum 2 characters is required"`
	Email    string `json:"email" binding:"required,email" message:"valid email is required"`
	Password string `json:"password" binding:"required,min=6" message:"password with minimum 6 characters is required"`
}
