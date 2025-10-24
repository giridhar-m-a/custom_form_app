package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

// GoogleAuthHandler handles Google OAuth authentication
// @Summary      Initiate Google OAuth authentication
// @Description  Redirects user to Google OAuth consent screen for authentication
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        code  query  string  true  "Authorization code from Google OAuth"
// @Success      200    {object}  object{status=int,message=string,data=object{accessToken=string,refreshToken=string,user=object{id=string,email=string,fullName=string,profilePic=string,profilePicId=string,createdAt=string,updatedAt=string}}}  "Authentication successful"
// @Failure      400    {object}  object{status=int,message=string}  "Bad request"
// @Failure      401    {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500    {object}  object{status=int,message=string}  "Internal server error"
// @Router       /auth/google [get]
// GoogleAuth registers the GET /auth/google endpoint on rg to initiate Google OAuth authentication.
// rg is the gin router group used to register the route.
func GoogleAuth(rg *gin.RouterGroup) {
	rg.GET("/auth/google", handler.GoogleAuthHandler)
}