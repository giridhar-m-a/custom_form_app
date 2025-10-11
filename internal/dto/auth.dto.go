package dto

import serializers "github.com/giridhar-m-a/custom_form_app/internal/serialisers"



type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         serializers.User   `json:"user"`
}

type GoogleAuthRequest struct {
	Code string `json:"code" binding:"required" form:"code" message:"code is required"`
}