package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Check API health
// @Description Returns OK if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func RegisterHealth(rg *gin.RouterGroup) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}

