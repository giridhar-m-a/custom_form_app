package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtSvc := services.NewJWTService()
		authHeader := c.GetHeader("Authorization")

		claims, err := jwtSvc.ValidateInvitationToken(authHeader)
		if err != nil {
			log.Printf("[Middleware] Authentication failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
			return
		}

		formRepo := repositories.NewFormsRepository(db.Queries)
		form, err := formRepo.GetFormByID(claims.FormID, c)
		if err != nil {
			log.Printf("[Middleware] Form not found: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
			return
		}

		if form.FormStatus.FormStatus != sqlc.FormStatusPublished {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Form is not published", "status": http.StatusUnauthorized})
			return
		}

		if form.FormAccess.FormAccess != sqlc.FormAccessPublic && claims.InvitationID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to access this form", "status": http.StatusUnauthorized})
			return
		}

		if form.ClosingTime.Valid && form.ClosingTime.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Form is closed", "status": http.StatusUnauthorized})
			return
		}

		// Store userID in context
		c.Set("invitationID", claims.InvitationID)
		c.Set("formID", claims.FormID)
		c.Next()
	}
}
