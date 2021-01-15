package user

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bri-bridge/backend-bridge-api/models"
	"gitlab.com/bri-bridge/backend-bridge-api/requests"
)

// IUserUsecase ..
type IUserUsecase interface {
	Register(c *gin.Context, request requests.RegisterRequest) (string, string, error)
	Login(c *gin.Context, request requests.LoginRequest) (*models.UserResponse, error)
}
