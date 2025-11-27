package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func Users(rg *gin.RouterGroup) {

	usersHandler := handler.NewUsersHandler()
	
	rg.GET("/users/me", usersHandler.GetMe)

}