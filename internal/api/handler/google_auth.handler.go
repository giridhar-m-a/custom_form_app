package handler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	serializers "github.com/giridhar-m-a/custom_form_app/internal/serialisers"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
)

// GoogleAuthHandler handles Google OAuth authentication by exchanging the provided
// authorization code for a user identity, issuing an access token and a refresh
// token, and returning a JSON payload containing the tokens and the serialized user.
//
// It reads token expiration values from the environment (JWT_EXPIRATION_HOURS and
// JWT_REFRESH_EXPIRATION_HOURS) with sensible defaults, and determines the token
// audience from the request Origin header or falls back to the request Host.
// On error it responds with appropriate HTTP status codes: 400 for invalid query
// parameters, 401 for authentication failure, and 500 for token generation errors.
func GoogleAuthHandler(c *gin.Context) {
	var query dto.GoogleAuthRequest
	// domain := os.Getenv("APP_DOMAIN")
	log.Printf("GoogleAuthHandler: Received query params: %+v", c.Request.URL.Query())
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Printf("Error binding query params: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	userRepo := repositories.NewSQLCUserRepository(db.Queries)
	userService := services.UserServiceProvider(userRepo)
	googleAuthService := services.NewGoogleAuthService(userService)
	authService := services.NewAuthService(googleAuthService)
	jwtService := services.NewJWTService()

	user, err := authService.AuthenticateWithGoogle(c.Request.Context(), query.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
		log.Printf("Error in GoogleAuthHandler: %v", err)
		return
	}

	jwtExpire, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		jwtExpire = 3600 // default to 1 hour if not set
	}
	refreshExpire, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRATION_HOURS"))
	if err != nil {
		refreshExpire = 86400 // default to 24 hours if not set
	}

	audience := c.Request.Header.Get("Origin")
	if audience == "" {
		audience = c.Request.Host // fallback if Origin is not set
	}

	token, err := jwtService.GenerateToken(user.UserID.String(), time.Duration(jwtExpire)*1*time.Second, audience)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "status": http.StatusInternalServerError})
		log.Printf("Error generating token: %v", err)
		return
	}
	refreshToken, err := jwtService.GenerateToken(user.UserID.String(), time.Duration(refreshExpire)*1*time.Second, audience)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "status": http.StatusInternalServerError})
		log.Printf("Error generating token: %v", err)
		return
	}

	// c.SetCookie("auth_token", token, 3600, "/", domain, false, true)
	// c.SetCookie("refresh_token", refreshToken, 86400, "/", domain, false, true)
	serializedUser := serializers.MapGetUserByGoogleIdRow(user)

	c.JSON(http.StatusOK, gin.H{"data": dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         serializedUser,
	}, "status": http.StatusOK, "message": "Authentication successful"})
}