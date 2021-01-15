package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler initialize health check handler
type HealthCheckHandler struct{}

// Check for database connection
func (a *HealthCheckHandler) Check(c *gin.Context) {
	// dbTimeStamp := a.HealthCheckUsecase.GetDBTimestamp()
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
	return
}
