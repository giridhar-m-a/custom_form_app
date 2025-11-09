package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func Auth(rg *gin.RouterGroup) {
	authHandler := handler.NewAuthHandler()

	rg.GET("/auth/google", authHandler.GoogleAuthHandler)

	rg.POST("/auth/register", authHandler.UserRegisterHandler)

	rg.POST("/auth/login", authHandler.EmailPasswordAuthHandler)
}
