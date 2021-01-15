package repository

import (
	"backend-sample-la/app/user"
	"backend-sample-la/models"
	"backend-sample-la/requests"

	gorm "github.com/jinzhu/gorm"
	// "gorm.io/gorm"
)

type userRepository struct {
	Conn *gorm.DB
}

// NewUserRepository ..
func NewUserRepository(Conn *gorm.DB) user.Repository {
	return &userRepository{Conn}
}

func (ur *userRepository) GetRoleID(name string) (int64, error) {

	var role models.Role

	if err := ur.Conn.Where("name = ?", name).First(&role).Error; err != nil {
		return 0, err
	}

	return role.Id, nil
}

func (ur *userRepository) GetUserById(userID string) (user models.User) {

	if err := ur.Conn.Model(user).Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}
	}
	return
}

func (ur *userRepository) Insert(userDoc *models.User) (*string, error) {

	if err := ur.Conn.Create(&userDoc).Error; err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur *userRepository) Get(email string) (*models.User, error) {
	var (
		user models.User
	)
	if err := ur.Conn.Where("email = ?", email).Find(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil

}

func (ur *userRepository) GetUsers() (*[]models.User, error) {
	var (
		users []models.User
	)
	if err := ur.Conn.Find(&users).Error; err != nil {
		return &users, nil
	}

	return &users, nil
}

func (ur *userRepository) DeleteUser(id string) (*models.User, error) {
	var (
		user models.User
	)
	if err := ur.Conn.Where("id = ?", id).Delete(&user).Error; err != nil {
		return &user, nil
	}

	return &user, nil
}

// UpdateUser ..
func (ur *userRepository) UpdateUser(req requests.UpdateUser) error {
	var (
		err  error
		user models.User
	)

	updateUser := map[string]interface{}{"Name": req.Name, "RoleId": req.RoleId}

	if err = ur.Conn.Model(&user).Where("id = ?", req.Id).Updates(updateUser).Error; err != nil {
		return nil
	}

	return err
}

// GetById ..
func (ur *userRepository) GetById(ID string) (*models.User, error) {
	var (
		user models.User
	)
	if err := ur.Conn.Where("id = ?", ID).Find(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}
