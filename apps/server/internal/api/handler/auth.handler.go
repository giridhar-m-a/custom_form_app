package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	serializers "github.com/giridhar-m-a/custom_form_app/internal/serialisers"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/google/uuid"
)

type AuthHandler interface {
	GoogleAuthHandler(c *gin.Context)
	UserRegisterHandler(c *gin.Context)
	EmailPasswordAuthHandler(c *gin.Context)
}

type authHandler struct {
	userService       services.UserService
	authService       services.AuthService
	googleAuthService services.GoogleAuthService
}

func NewAuthHandler() AuthHandler {
	userRepo := repositories.NewSQLCUserRepository(db.Queries)
	userService := services.UserServiceProvider(userRepo)
	googleAuthService := services.NewGoogleAuthService(userService)

	authService := services.NewAuthService(googleAuthService, userService)
	return &authHandler{
		userService:       userService,
		authService:       authService,
		googleAuthService: googleAuthService,
	}
}

// @Summary      Initiate Google OAuth authentication
// @Description  Redirects user to Google OAuth consent screen for authentication
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        code  query  string  true  "Authorization code from Google OAuth"
// @Success      200  {object}  object{status=int,message=string,data=object{accessToken=string,refreshToken=string,user=object{id=string,email=string,fullName=string,profilePic=string,profilePicId=string,createdAt=string,updatedAt=string}}}  "Authentication successful"
// @Failure      400  {object}  object{status=int,message=string}  "Bad request"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /auth/google [get]
// @Schemes      https
func (a *authHandler) GoogleAuthHandler(c *gin.Context) {
	var query dto.GoogleAuthRequest
	// domain := os.Getenv("APP_DOMAIN")
	log.Printf("GoogleAuthHandler: Received query params: %+v", c.Request.URL.Query())
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Printf("Error binding query params: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	user, err := a.authService.AuthenticateWithGoogle(c.Request.Context(), query.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
		log.Printf("Error in GoogleAuthHandler: %v", err)
		return
	}

	audience := c.Request.Header.Get("Origin")
	if audience == "" {
		audience = c.Request.Host // fallback if Origin is not set
	}

	token, refreshToken, err := a.authService.GenerateTokens(user.UserID.String(), audience)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate tokens", "status": http.StatusInternalServerError})
		log.Printf("Error generating tokens: %v", err)
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

// @Summary      Login with email and password
// @Description  Authenticates a user and returns access & refresh tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body  dto.EmailPasswordAuthRequest  true  "Login credentials"
// @Success      200  {object}  object{status=int,message=string,data=object{accessToken=string,refreshToken=string,user=object{id=string,email=string,fullName=string,profilePic=string,profilePicId=string,createdAt=string,updatedAt=string}}}  "Login successful"
// @Failure      400  {object}  object{status=int,message=string}  "Invalid credentials"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /auth/login [post]
// @Schemes      https
func (a *authHandler) EmailPasswordAuthHandler(c *gin.Context) {
	var payload dto.EmailPasswordAuthRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	user, err := a.authService.AuthenticateWithEmailPassword(c.Request.Context(), payload)
	if err != nil {
		fmt.Printf("Error Getting User %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
		return
	}

	audience := c.Request.Header.Get("Origin")
	if audience == "" {
		audience = c.Request.Host // fallback if Origin is not set
	}

	token, refreshToken, err := a.authService.GenerateTokens(user.UserID.String(), audience)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate tokens", "status": http.StatusInternalServerError})
		log.Printf("Error generating tokens: %v", err)
		return
	}

	serializedUser := serializers.MapGetUserByEmailRow(user)

	c.JSON(http.StatusOK, gin.H{"data": dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         serializedUser,
	}, "status": http.StatusOK, "message": "Authentication successful"})
}

// @Summary      Register a new user
// @Description  Registers a user with name, email, and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body  dto.EmailPasswordRegisterRequest  true  "User registration payload"
// @Success      200  {object}  object{status=int,message=string,data=object{accessToken=string,refreshToken=string,user=object{id=string,email=string,fullName=string,profilePic=string,profilePicId=string,createdAt=string,updatedAt=string}}}  "Authentication successful"
// @Failure      400  {object}  object{status=int,message=string}  "Invalid request"
// @Failure      409  {object}  object{status=int,message=string}  "Email already exists"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /auth/register [post]
// @Schemes      https
func (a *authHandler) UserRegisterHandler(c *gin.Context) {
	var payload dto.EmailPasswordRegisterRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	existingUser, err := a.userService.GetUserDetailsByEmail(c.Request.Context(), payload.Email)
	if err == nil && existingUser.UserID != uuid.Nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Something went wrong", "status": http.StatusConflict})
		return
	}

	user, err := a.authService.CreateUserWithEmailPassword(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	token, refreshToken, err := a.authService.GenerateTokens(user.UserID.String(), c.Request.Host)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate tokens", "status": http.StatusInternalServerError})
		log.Printf("Error generating tokens: %v", err)
		return
	}

	serializedUser := serializers.MapUser(user)

	c.JSON(http.StatusOK, gin.H{"data": dto.AuthResponse{User: serializedUser, AccessToken: token, RefreshToken: refreshToken}, "status": http.StatusOK, "message": "User registered successfully"})
}
