package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	queryparams "github.com/giridhar-m-a/custom_form_app/internal/repositories/query-params"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

func GoogleAuthHandler(c *gin.Context) {
	var query queryparams.GoogleAuthQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
        return
    }
    token, err := services.GoogleAuthService(query.Code)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
        log.Printf("Error in GoogleAuthHandler: %v", err)
        return
    }
    googleUserInfo, err := utils.GetUserDetails(token.AccessToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "status": http.StatusUnauthorized})
        log.Printf("Error fetching user info: %v", err)
        return
    }
    existingUser, err := db.Queries.GetUserByGoogleId(db.Ctx, googleUserInfo["id"].(string))
    if err == nil && existingUser.UserID != uuid.Nil {
        log.Printf("User already exists: %v", existingUser.UserFullName)
        c.JSON(http.StatusOK, gin.H{"user": existingUser, "status": http.StatusOK, "message": "Authentication successful"})
        return
    }
    userParams := sqlc.CreateUserParams{
        UserFullName:     googleUserInfo["name"].(string),
        UserEmail:        googleUserInfo["email"].(string),
        UserGoogleID:     googleUserInfo["id"].(string),
        UserProfilePicID: uuid.NullUUID{}, // Initialize as NullUUID
    }

    user, err := db.Queries.CreateUser(db.Ctx, userParams)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
        log.Printf("Error creating user: %v", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"user": user, "status": http.StatusOK, "message": "Authentication successful"})
}