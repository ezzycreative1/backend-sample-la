package app

import (
	"github.com/gin-gonic/gin"
	HCHandler "backend-sample-la/app/healthcheck/handler"
)

// HealthCheckHTTPHandler ..
func HealthCheckHTTPHandler(router *gin.Engine) {
	handler := &HCHandler.HealthCheckHandler{}
	router.GET("/check", handler.Check)
}
