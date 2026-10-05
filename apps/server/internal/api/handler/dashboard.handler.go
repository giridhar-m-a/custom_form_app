package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)



type DashboardHandler interface {
	GetDashboardData(ctx *gin.Context)
}

type dashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler() DashboardHandler {
	return &dashboardHandler{
		dashboardService: services.NewDashboardService(repositories.NewDashboardRepository(db.Queries)),
	}
}

// @Summary      Get dashboard data
// @Description  Get dashboard data for the authenticated user
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Success      200   {object}  dto.ApiResponse[dto.DashboardResponse]  "Dashboard data fetched successfully"
// @Failure      400   {object}  dto.ApiResponse[any]  "Invalid request payload"
// @Failure      401   {object}  dto.ApiResponse[any]  "Unauthorized"
// @Failure      500   {object}  dto.ApiResponse[any]  "Internal server error"
// @Router       /dashboard [get]
// @Security     BearerAuth
func (h *dashboardHandler) GetDashboardData(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}

	nulUUid:= utils.ConvertStringToNullUUID(userID.(string))
	dashboardData, err := h.dashboardService.GetDashboardData(ctx, nulUUid)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	response:= dto.ApiResponse[dto.DashboardResponse]{
		Status: 200,
		Message: "Dashboard data fetched successfully",
		Data: dashboardData,
	}

	ctx.JSON(200, response)
}