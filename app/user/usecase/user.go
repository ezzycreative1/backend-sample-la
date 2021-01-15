package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"backend-sample-la/app/user"
	userInterfaces "backend-sample-la/app/user"
	"backend-sample-la/helpers"
	"backend-sample-la/models"
	"backend-sample-la/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userUsecase struct {
	repo user.Repository
}

type M map[string]interface{}

// NewUserUsecase ..
func NewUserUsecase(repo user.Repository) userInterfaces.IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

// Register ..
func (us *userUsecase) Register(c *gin.Context, request requests.RegisterRequest) (string, string, error) {

	registerDoc := models.User{}

	user, err := us.repo.Get(request.Email)

	if user.Email != "" {

		if user.VerifedAt == "" {
			return user.Id, user.VerifyToken, nil
		}

		return "", "", errors.New("Email In Use")
	}
	// Hash Password
	password, _ := hashPassword(request.Password)

	// Generate Random uuID
	u, err := uuid.NewRandom()
	if err != nil {
		log.Println(fmt.Errorf("Failed to generate UUID: %w", err))
		return "", "", errors.New("Failed to generate UUID")

	}

	now := time.Now()

	registerDoc.Id = u.String()
	registerDoc.Password = password
	registerDoc.Email = request.Email
	registerDoc.Name = request.Name
	registerDoc.RoleId = 3 // user default
	registerDoc.VerifyToken = helpers.GenerateRandString(18)

	registerDoc.CreatedAt = now
	registerDoc.UpdatedAt = now
	// Save
	_, err = us.repo.Insert(&registerDoc)

	if err != nil {
		return "", "", err
	}

	var key = fmt.Sprintf("%s-%s", registerDoc.Id, registerDoc.VerifyToken)
	cache.Set(key, request.Email, time.Duration(15*time.Minute))
	// token, err := GenerateToken(request.Email)
	return registerDoc.Id, registerDoc.VerifyToken, nil
}

// Login ..
func (us *userUsecase) Login(c *gin.Context, request requests.LoginRequest) (*models.UserResponse, error) {
	res := &models.UserResponse{}
	if request.Email == "" || request.Password == "" {
		return nil, errors.New("Email Or Password Cannot Blank")
	}

	checkEmail := helpers.ValidEmail(request.Email)
	if !checkEmail {
		return nil, errors.New("Format Email Wrong")
	}

	user, _ := us.repo.Get(request.Email)

	if user.VerifedAt == "" {
		return nil, errors.New("Unverified account. Please check your email to verify")
	}
	// fmt.Println(user.)

	checkPassword := checkPasswordHash(request.Password, user.Password)

	if !checkPassword {

		type FailCount struct {
			Value int `json:"value"`
		}

		// check key exist
		var key = fmt.Sprintf("fail-%s", request.Email)
		if cache.Exists(key) {

			var value FailCount

			val, err := cache.Get(key)
			if err != nil {
				fmt.Println(err.Error())
			} else {

				go cache.Del(key)
				if err := json.Unmarshal([]byte(val), &value); err != nil {
					fmt.Println(err.Error())
				} else {

					if value.Value == 2 {
						var newKey = fmt.Sprintf("susspend-%s", request.Email)
						go cache.Set(newKey, "suspend", time.Duration(30*time.Minute))

					} else {
						value.Value += 1
						jsByte, _ := json.Marshal(value)
						go cache.Set(key, string(jsByte), time.Duration(15*time.Minute))
					}
				}
			}
		} else {

			var valStr = FailCount{
				1,
			}
			jsByte, _ := json.Marshal(valStr)
			var val = string(jsByte)

			err := cache.Set(key, val, time.Duration(15*time.Minute))
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		return nil, errors.New("Email Or Password Fail")
	}

	role := getRoleName(user.RoleId)

	getToken, err := GenerateToken(user.Email, user.Id, role)

	if err != nil {
		return nil, errors.New("Genarete Token Failed")
	}

	res.Role = getRoleName(user.RoleId)
	res.Token = getToken

	return res, nil
}
