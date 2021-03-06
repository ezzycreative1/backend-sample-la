package user

import (
	"backend-sample-la/models"
	"backend-sample-la/requests"
)

// Repository ..
type Repository interface {
	Insert(*models.User) (*string, error)
	Get(string) (*models.User, error)
	GetById(string) (*models.User, error)
	GetUsers() (*[]models.User, error)
	DeleteUser(string) (*models.User, error)
	UpdateUser(req requests.UpdateUser) error
	GetRoleID(name string) (int64, error)
	GetUserById(userID string) (user models.User)
}
