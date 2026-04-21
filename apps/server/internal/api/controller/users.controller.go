package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func Users(rg *gin.RouterGroup) {

	usersHandler := handler.NewUsersHandler()

	api := rg.Group("/users/me")

	api.GET("", usersHandler.GetMe)
	api.PUT("/profile-pic", usersHandler.UpdateProfilePic)
	api.PATCH("", usersHandler.UpdateUser)
	api.PUT("/update-password", usersHandler.UpdatePassword)
	api.DELETE("", usersHandler.DeleteUser)

}
