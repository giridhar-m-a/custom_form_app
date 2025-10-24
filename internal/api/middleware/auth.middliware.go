package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
)

// AuthMiddleware returns a Gin middleware handler that validates a JWT from the request's
// Authorization header and, on success, stores the extracted user ID in the request context.
// On validation failure it aborts the request with HTTP 401 and a JSON error message (context key "userID").
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtSvc := services.NewJWTService()
		authHeader := c.GetHeader("Authorization")
		userID, err := jwtSvc.ValidateToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// Store userID in context
		c.Set("userID", userID)
		c.Next()
	}
}