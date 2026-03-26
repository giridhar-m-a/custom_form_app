package handler

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/cache"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type UsersHandler interface {
	GetMe(ctx *gin.Context)
	UpdateProfilePic(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	DeleteUserProfilePic(ctx *gin.Context)
}

type usersHandler struct {
	userService services.UserService
	bucket      string
}

func NewUsersHandler() UsersHandler {
	userRepo := repositories.NewSQLCUserRepository(db.Queries)
	userService := services.UserServiceProvider(userRepo)
	bucket := utils.GetEnv("MINIO_BUCKET_NAME", "custom-form-app")
	return &usersHandler{userService: userService, bucket: bucket}
}

// @Summary      Get current user details
// @Description  Retrieves details of the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse[dto.User]  "User retrieved successfully"
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
	var response dto.ApiResponse[dto.User]
	key := "user:userID:" + userID.(string)
	cachedUser, err := cache.Get(ctx, key)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &response); err == nil {
			// ✅ Cache hit, return immediately
			ctx.JSON(200, response)
			return
		}
	}
	user, err := h.userService.GetUserDetailsById(ctx, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	// userResponse := serializers.MapGetUserByIDRow(user)
	var profilepic string

	if user.FileName.Valid {
		// signed, err := services.GetMinioSignedURL(h.bucket, user.FileName.String, time.Hour*24, "")
		if err == nil {
			profilepic = user.FileName.String
		}

	}

	response = dto.ApiResponse[dto.User]{
		Status:  200,
		Message: "User retrieved successfully",
		Data: dto.User{
			UserID:         user.UserID.String(),
			UserEmail:      user.UserEmail,
			UserFullName:   user.UserFullName,
			UserCreatedAt:  user.UserCreatedAt.Time,
			UserUpdatedAt:  user.UserUpdatedAt.Time,
			UserProfilePic: profilepic,
		},
	}

	userJSON, _ := json.Marshal(response)
	_ = cache.Set(ctx, key, string(userJSON))
	ctx.JSON(200, response)
}

// @Summary      Update user details
// @Description  Updates the details of the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        form  body      dto.UpdateUserDetailsDTO true  "Form data"
// @Success      200  {object}  dto.ApiResponse[dto.User]  "User updated successfully"
// @Failure      400  {object}  object{status=int,message=string}  "Bad request"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /users/me [patch]
// @Schemes      https
// @Security     BearerAuth
func (h *usersHandler) UpdateUser(ctx *gin.Context) {

	var data dto.UpdateUserDetailsDTO

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}

	user, err := h.userService.UpdateUser(ctx, userID.(string), dto.UserUpdateDTO{
		UserFullName: data.UserFullName,
	})
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, dto.ApiResponse[dto.User]{
		Status:  200,
		Message: "User retrieved successfully",
		Data: dto.User{
			UserID:        user.UserID.String(),
			UserEmail:     user.UserEmail,
			UserFullName:  user.UserFullName,
			UserCreatedAt: user.UserCreatedAt.Time,
			UserUpdatedAt: user.UserUpdatedAt.Time,
		},
	})
}

// @Summary      Change password
// @Description  Changes the password of the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        form  body      dto.UpdateUserPasswordDTO true  "Form data"
// @Success      200  {object}  dto.ApiResponse[dto.User]  "Updated User Password successfully"
// @Failure      400  {object}  object{status=int,message=string}  "Bad request"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /users/me/update-password [put]
// @Schemes      https
// @Security     BearerAuth
func (h *usersHandler) UpdatePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}

	var data dto.UpdateUserPasswordDTO

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	if data.UserPassword != data.UserVerifyPassword {
		utils.HandleError(ctx, errors.New("passwords do not match"))
		return
	}

	oldHashed, err := h.userService.GetUserPassword(ctx, userID.(string))

	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	if oldHashed == "" {
		utils.HandleError(ctx, errors.New("User not created with credentials"))
		return
	}

	bcrypt := services.NewBcryptService()

	isValid := bcrypt.ComparePassword(oldHashed, data.OldPassword)

	if !isValid {
		utils.HandleError(ctx, errors.New("Invalid password"))
		return
	}

	isValid = bcrypt.ComparePassword(oldHashed, data.UserPassword)

	if isValid {
		utils.HandleError(ctx, errors.New("Don't Use Same Password"))
		return
	}

	newHashed, err := bcrypt.HashPassword(data.UserPassword)

	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	user, err := h.userService.UpdateUser(ctx, userID.(string), dto.UserUpdateDTO{
		UserPassword: newHashed,
	})
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, dto.ApiResponse[dto.User]{
		Status:  200,
		Message: "User password updated successfully",
		Data: dto.User{
			UserID:        user.UserID.String(),
			UserEmail:     user.UserEmail,
			UserFullName:  user.UserEmail,
			UserCreatedAt: user.UserCreatedAt.Time,
			UserUpdatedAt: user.UserUpdatedAt.Time,
		},
	})
}

// @Summary      Delete User
// @Description  Deletes the authenticated user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse[any]  "User deleted successfully"
// @Failure      400  {object}  object{status=int,message=string}  "Bad request"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /users/me [delete]
// @Schemes      https
// @Security     BearerAuth
func (h *usersHandler) DeleteUser(ctx *gin.Context) {

	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}

	err := h.userService.DeleteUser(ctx, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "User deleted successfully",
	})
}

// @Summary      Delete User Profile Picture
// @Description  Deletes the authenticated user's profile picture
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse[any]  "User profile picture deleted successfully"
// @Failure      400  {object}  object{status=int,message=string}  "Bad request"
// @Failure      401  {object}  object{status=int,message=string}  "Unauthorized"
// @Failure      500  {object}  object{status=int,message=string}  "Internal server error"
// @Router       /users/me/profile-pic [delete]
// @Schemes      https
// @Security     BearerAuth
func (h *usersHandler) DeleteUserProfilePic(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}

	err := h.userService.DeleteUserProfilePic(ctx, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.JSON(200, dto.ApiResponse[any]{
		Status:  200,
		Message: "Profile picture deleted successfully",
	})
}

// @Summary Upload a User Profile Picture
// @Description Update the user's profile picture
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} dto.ApiResponse[dto.UserProfilePicResponse] "Profile picture updated successfully"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/profile-pic [put]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *usersHandler) UpdateProfilePic(ctx *gin.Context) {
	// Parse the multipart form, with a max memory of 10MB
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Retrieve the file from form data
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	defer file.Close()

	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}
	key := "user:userID:" + userID.(string)
	profile, err := h.userService.UpdateUserProfilePic(ctx, userID.(string), dto.FileUploadPayload{File: file, FileInfo: header})
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	userResponse := dto.UserProfilePicResponse{
		UserID:   profile.UserID.String(),
		FileName: profile.FileName,
		FileSize: profile.FileSize,
		FileType: profile.FileType,
	}
	ctx.JSON(200, dto.ApiResponse[dto.UserProfilePicResponse]{
		Status:  200,
		Message: "Profile picture updated successfully",
		Data:    userResponse,
	})
	cache.Del(ctx, key)
}
