package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	serializers "github.com/giridhar-m-a/custom_form_app/internal/serialisers"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type UsersHandler interface {
	GetMe(ctx *gin.Context)
}

type usersHandler struct {
	userService services.UserService
}

func NewUsersHandler() UsersHandler {
	userRepo := repositories.NewSQLCUserRepository(db.Queries)
	userService := services.UserServiceProvider(userRepo)
	return &usersHandler{userService: userService}
}

// @Summary      Get current user details
// @Description  Retrieves details of the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  object{status=int,message=string,data=object{id=string,email=string,fullName=string,profilePic=string,profilePicId=string,createdAt=string,updatedAt=string}}  "User retrieved successfully"
// @Failure      400  {object}  object{status=int,message=string}  "Bad request"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /users/me [get]
// @Schemes      https
// @Security     BearerAuth
func (h *usersHandler) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}
	user, err := h.userService.GetUserDetailsById(ctx, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	userResponse := serializers.MapGetUserByIDRow(user)
	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "User retrieved successfully",
		"data":    userResponse,
	})
}
