package user

import (
	"backend-sample-la/models"
	"backend-sample-la/requests"

	"github.com/gin-gonic/gin"
)

// IUserUsecase ..
type IUserUsecase interface {
	Register(c *gin.Context, request requests.RegisterRequest) (string, string, error)
	Login(c *gin.Context, request requests.LoginRequest) (*models.UserResponse, error)
}
