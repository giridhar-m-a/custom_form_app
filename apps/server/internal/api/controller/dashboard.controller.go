package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)



func  GetDashboardData(rg *gin.RouterGroup) {
	dashboardHandler := handler.NewDashboardHandler()
	api:= rg.Group("/dashboard")

	api.GET("/", dashboardHandler.GetDashboardData)
}
