package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func Auth(rg *gin.RouterGroup) {
	authHandler := handler.NewAuthHandler()

	auth := rg.Group("/auth")

	auth.GET("/google", authHandler.GoogleAuthHandler)

	auth.POST("/register", authHandler.UserRegisterHandler)

	auth.POST("/login", authHandler.EmailPasswordAuthHandler)

	auth.GET("/verify", authHandler.VerifyTokenHandler)

	auth.GET("/refresh-token", authHandler.RefreshTokenHandler)
}
