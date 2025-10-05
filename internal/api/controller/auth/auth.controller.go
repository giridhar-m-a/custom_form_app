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
// @Success      200   {object}  map[string]interface{}  "Authentication successful"
// @Router       /auth/google [get]
// @Schemes http
func GoogleAuth(rg *gin.RouterGroup) {
	rg.GET("/auth/google", handler.GoogleAuthHandler)
}
