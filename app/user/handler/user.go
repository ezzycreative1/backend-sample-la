package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"backend-sample-la/cache"

	// "net/smtp"

	BaseHandler "backend-sample-la/app/base/handler"
	UsInterfaces "backend-sample-la/app/user"
	"backend-sample-la/helpers/blog"
	res "backend-sample-la/models"
	"backend-sample-la/requests"

	"github.com/gin-gonic/gin"
)

// UserHandler ..
type UserHandler struct {
	RUsecase UsInterfaces.IUserUsecase
}

// Register ..
func (u *UserHandler) Register(c *gin.Context) {
	requestBody := requests.RegisterRequest{}
	err := c.BindJSON(&requestBody)
	if err != nil {

		// save to log
		go blog.LevelError("Register", "/user/register", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	// this variabel for logging
	reqJSON, _ := json.Marshal(requestBody)

	// validate request
	err = u.RUsecase.RequestRegisterValid(requestBody)
	if err != nil {
		go blog.LevelDebug("Register", "/user/register", string(reqJSON), err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	//userID, token,
	userID, token, err := u.RUsecase.Register(c, requestBody)
	fmt.Println(token)
	data := res.RegisterResponse{
		Message: "please check your email to validate your account",
	}

	if err != nil {

		// save to log
		go blog.LevelError("Register", "/user/register", string(reqJSON), err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	go u.RUsecase.SendValidationEmail(requestBody.Name, requestBody.Email, userID, token)

	// save to log
	dataJSON, _ := json.Marshal(data)
	go blog.LevelInfo("Register", "/user/register", string(reqJSON), string(dataJSON), "", time.Now())

	BaseHandler.RespondSuccess(c, "", data)
	return
}

// Login ..
func (u *UserHandler) Login(c *gin.Context) {
	requestBody := requests.LoginRequest{}
	err := c.BindJSON(&requestBody)
	if err != nil {

		// save to log
		go blog.LevelError("Login", "/user/login", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	// validate
	if err := u.RUsecase.RequestLoginValid(requestBody); err != nil {
		// save to log
		go blog.LevelDebug("Login", "/user/login", "", err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	var key = fmt.Sprintf("susspend-%s", requestBody.Email)
	fmt.Println(cache.Exists(key))
	if cache.Exists(key) {
		BaseHandler.FailedResponseBackend(c, fmt.Errorf("your account has been suspended 30 menute"))
		return
	}

	fmt.Println(requestBody.Password)
	//return
	resLogin, err := u.RUsecase.Login(c, requestBody)
	reqJSON, _ := json.Marshal(requestBody)

	if err != nil {

		go blog.LevelError("Login", "/user/login", string(reqJSON), err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	resJSON, _ := json.Marshal(resLogin)
	go blog.LevelInfo("Login", "/user/login", string(reqJSON), string(resJSON), "", time.Now())

	BaseHandler.RespondSuccess(c, "", resLogin)
	return
}
