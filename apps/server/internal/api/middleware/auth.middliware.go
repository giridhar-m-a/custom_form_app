package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtSvc := services.NewJWTService()
		authHeader := c.GetHeader("Authorization")

		userID, err := jwtSvc.ValidateToken(authHeader)
		if err != nil {
			log.Printf("[Middleware] Authentication failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
			return
		}
		// Store userID in context
		c.Set("userID", userID)
		c.Next()
	}
}
