package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/bri-bridge/backend-bridge-api/cache"

	// "net/smtp"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gin-gonic/gin"
	BaseHandler "gitlab.com/bri-bridge/backend-bridge-api/app/base/handler"
	UsInterfaces "gitlab.com/bri-bridge/backend-bridge-api/app/user"
	middleware "gitlab.com/bri-bridge/backend-bridge-api/app/user/usecase"
	"gitlab.com/bri-bridge/backend-bridge-api/helpers/blog"
	"gitlab.com/bri-bridge/backend-bridge-api/models"
	res "gitlab.com/bri-bridge/backend-bridge-api/models"
	"gitlab.com/bri-bridge/backend-bridge-api/requests"
)

// UserHandler ..
type UserHandler struct {
	RUsecase UsInterfaces.IUserUsecase
}

// ValidateGoogleToken ..
func (u *UserHandler) ValidateGoogleToken(c *gin.Context) {

	var req requests.ValidateGoogleToken

	if err := c.ShouldBindJSON(&req); err != nil {

		// save to log
		go blog.LevelError("ValidateGoogleToken", "/user/google/account", "", err.Error(), "", time.Now())

		e := errors.New("user not found, or invalid email")

		c.Error(e)
		BaseHandler.FailedResponseBackend(c, e)
		return
	}

	v := validator.New()
	if err := v.Struct(req); err != nil {
		// save activity log
		go blog.LevelError("ValidateGoogleToken", "/user/google/account", "", err.Error(), "", time.Now())

		e := errors.New("user not found, or invalid email")

		c.Error(e)
		BaseHandler.FailedResponseBackend(c, e)
		return
	}

	reqJSON, _ := json.Marshal(req)

	// get user info from google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + req.Token)
	if err != nil {
		// save activity log
		go blog.LevelError("ValidateGoogleToken", "/user/google/account", string(reqJSON), err.Error(), "", time.Now())

		e := errors.New("user not found, or invalid email")

		c.Error(e)
		BaseHandler.FailedResponseBackend(c, e)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// save activity log
		go blog.LevelError("ValidateGoogleToken", "/user/google/account", string(reqJSON), "", "", time.Now())

		e := errors.New("user not found, or invalid email")
		fmt.Println(resp.StatusCode)

		c.Error(e)
		BaseHandler.FailedResponseBackend(c, e)
		return
	}

	// read response body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// save activity log
		go blog.LevelError("ValidateGoogleToken", "/user/google/account", string(reqJSON), err.Error(), "", time.Now())

		e := errors.New("user not found, or invalid email")

		c.Error(e)
		BaseHandler.FailedResponseBackend(c, e)
		return
	}

	// decode json resp google to map[string]string
	var respGoogle map[string]string

	json.Unmarshal(content, &respGoogle)

	// check user with email exist or not and vaidate email has been connected with google
	res, err := u.RUsecase.RegisterWithGoogleAccount(respGoogle)
	if err != nil {
		// save activity log
		go blog.LevelError("ValidateGoogleToken", "/user/google/account", string(reqJSON), err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	resJSON, _ := json.Marshal(res)
	go blog.LevelInfo("ValidateGoogleToken", "/user/google/account", string(reqJSON), string(resJSON), "", time.Now())

	BaseHandler.RespondSuccess(c, "", res)
	return

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

func (u *UserHandler) SendEmailRegister(c *gin.Context) {
	type Request struct {
		Id   string `json:"id"`
		Code string `json:"code"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		// save to log
		go blog.LevelError("Register", "/user/send/mail", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, fmt.Errorf("cannot read request body"))
		return
	}
	req.Id = strings.TrimSpace(req.Id)
	req.Code = strings.TrimSpace(req.Code)

	if req.Id == "" {
		BaseHandler.FailedResponseBackend(c, fmt.Errorf("id required"))
		return
	}
	if req.Code == "" {
		BaseHandler.FailedResponseBackend(c, fmt.Errorf("code required"))
		return
	}
	if err := u.RUsecase.SendMoreMailRegister(req.Id, req.Code); err != nil {
		// save to log
		go blog.LevelError("Register", "/user/send/mail", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, fmt.Errorf("cannot send email"))
		return
	}

	data := res.RegisterResponse{
		Message: "please check your email to validate your account",
	}

	BaseHandler.RespondSuccess(c, "", data)
	return
}

//reset password
func (u *UserHandler) EmailResetPassword(c *gin.Context) {

	type Request struct {
		Email string `json:"email"`
	}

	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		// save to log
		e := fmt.Errorf("cannot read request body")
		go blog.LevelError("Login", "/user/login", "", e.Error(), "", time.Now())

		c.Error(e)
		BaseHandler.FailedResponseBackend(c, e)
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	v := validator.New()
	if err := v.Var(req.Email, "required,email"); err != nil {
		go blog.LevelWarn("Login", "/user/login", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	// check email exist on database
	exist := u.RUsecase.CheckEmailExistResetPassword(req.Email)
	if !exist {
		var err = fmt.Errorf("email not found")
		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return

	}

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

// FacebookLogin ..
func (u *UserHandler) FacebookLogin(c *gin.Context) {

	u.RUsecase.FacebookLogin(c)

	return

}

// FacebookCallback  ...
func (u *UserHandler) FacebookCallback(c *gin.Context) {
	model, err := u.RUsecase.FacebookCallback(c)

	if err != nil {
		return
	}

	res := models.AccessToken{}

	res.Id = model.Id
	res.Name = model.Name
	res.Token = model.Token

	BaseHandler.RespondSuccess(c, "", res)

	return
}

func (u *UserHandler) GetProfile(c *gin.Context) {

	userID, _ := middleware.GetRoleAndUserID(c)
	fmt.Println("halllo")

	user := u.RUsecase.GetUserProfile(userID)

	BaseHandler.RespondSuccess(c, "", user)
	return
}

// ForgotPassword ..
func (u *UserHandler) ForgotPassword(c *gin.Context) {

	// validate requiest

	var req requests.ForgotPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		// handle error
		return
	}

	v := validator.New()
	if err := v.Struct(req); err != nil {
		// handle error
		return
	}

	// check user exist
	if err := u.RUsecase.ForgotPassword(req.Email); err != nil {
		// handle error
	}

	BaseHandler.RespondSuccess(c, "", "please check your email to validate your account")
	return
}

// VerifyAccount is method for validate user email
// with userid and token in url param
// @param *gin.Context
func (u *UserHandler) VerifyAccount(c *gin.Context) {

	token := c.Param("verifyToken")
	userID := c.Param("userid")

	fmt.Println(token, userID)

	v := validator.New()
	if err := v.Var(token, "required"); err != nil {
		BaseHandler.FailedAPIValidation("token must be filled", c)
		return
	}
	if err := v.Var(userID, "required"); err != nil {
		BaseHandler.FailedAPIValidation("token must be filled", c)
		return
	}

	ok, err := u.RUsecase.VerifyAccount(token, userID)
	if err != nil {
		// handler error
		BaseHandler.RespondNotFound(c, "", "Not found", "user not found", "")
		return
	}
	if !ok {
		var location = fmt.Sprintf("https://dev-bridge.tcd-dev.id/#/register-failed?id=%s&code=%s",
			userID, token)
		c.Redirect(http.StatusPermanentRedirect, location)
		c.AbortWithStatus(http.StatusPermanentRedirect)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "https://dev-bridge.tcd-dev.id/#/register-success")
	c.AbortWithStatus(http.StatusPermanentRedirect)
	return

}

// ChangePhoto ..
func (u *UserHandler) ChangePhoto(c *gin.Context) {
	var err error

	userID, _ := middleware.GetRoleAndUserID(c)

	_, err = u.RUsecase.ChangePhoto(userID, c)
	if err != nil {
		// save to log
		go blog.LevelError("ChangePhoto", "/user/change-photo", "", err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	// save to log
	go blog.LevelInfo("ChangePhoto", "/user/change-photo", "", "Data Success Updated", "", time.Now())

	BaseHandler.RespondSuccess(c, "", "Data Success Updated")
	return
}

// ChangeName ..
func (u *UserHandler) ChangeName(c *gin.Context) {
	requestBody := requests.ChangeName{}
	err := c.BindJSON(&requestBody)
	if err != nil {

		// save to log
		go blog.LevelError("ChangeName", "/user/change-name", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	reqJSON, _ := json.Marshal(requestBody)

	_, err = u.RUsecase.ChangeName(c, requestBody)
	if err != nil {

		// save to log
		go blog.LevelError("ChangeName", "/user/change-name", string(reqJSON), err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	// save to log
	go blog.LevelInfo("ChangeName", "/user/change-name", string(reqJSON), "Data Success Updated", "", time.Now())

	BaseHandler.RespondSuccess(c, "", "Data Success Updated")
	return
}

// ChangePassword ..
func (u *UserHandler) ChangePassword(c *gin.Context) {

	requestBody := requests.ChangePassword{}

	err := c.BindJSON(&requestBody)
	if err != nil {

		// save to log
		go blog.LevelError("ChangePassword", "/user/change-password", "", err.Error(), "", time.Now())

		c.Error(err)
		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	err = u.RUsecase.ChangePassword(c, requestBody)
	reqJSON, _ := json.Marshal(requestBody)

	if err != nil {

		// save to log
		go blog.LevelError("ChangePassword", "/user/change-password", string(reqJSON), err.Error(), "", time.Now())

		BaseHandler.FailedResponseBackend(c, err)
		return
	}

	// save to log
	go blog.LevelError("ChangePassword", "/user/change-password", string(reqJSON), "Success Updated Password", "", time.Now())

	BaseHandler.RespondSuccess(c, "", "Success Updated Password")
	return
}

// this method for handling image in email
func (u *UserHandler) ImageEmail(c *gin.Context) {

	content := c.Param("content")

	if content == "logo" {
		c.Header("Content-type", "image/png")
		var dir = fmt.Sprintf("%s/src/gitlab.com/bri-bridge/backend-bridge-api/storage/email/logo.png", os.Getenv("GOPATH"))
		c.File(dir)
	} else if content == "email" {
		c.Header("Content-type", "image/png")
		var dir = fmt.Sprintf("%s/src/gitlab.com/bri-bridge/backend-bridge-api/storage/email/email.png", os.Getenv("GOPATH"))
		c.File(dir)
	} else if content == "pass" {
		c.Header("Content-type", "image/png")
		var dir = fmt.Sprintf("%s/src/gitlab.com/bri-bridge/backend-bridge-api/storage/email/pass.png", os.Getenv("GOPATH"))
		c.File(dir)
	} else if content == "adminverify" {
		c.Header("Content-type", "image/png")
		var dir = fmt.Sprintf("%s/src/gitlab.com/bri-bridge/backend-bridge-api/storage/email/admin_verif.png", os.Getenv("GOPATH"))
		c.File(dir)
	} else {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusNotFound)
	}

}

func (u *UserHandler) Dashboard(c *gin.Context) {

	userID, role := middleware.GetRoleAndUserID(c)

	if role == "User" {
		var data = u.RUsecase.DashboardUser(userID)
		BaseHandler.RespondSuccess(c, "", data)
		return
	} else if role == "Admin" || role == "Super Admin" || role == "Viewer" {
		var data = u.RUsecase.DashboardAdministration(userID)
		BaseHandler.RespondSuccess(c, "", data)
		return
	} else {
		BaseHandler.RespondForbidden(c)
		return
	}
}

func (u *UserHandler) ImageUser(c *gin.Context) {

	content := c.Param("content")
	if content == "" {
		BaseHandler.RespondNotFound(c, "04", "not found", "image not found", "")
		return
	}

	// get image
	image := u.RUsecase.GetImageUser(content)
	if image == "" {
		BaseHandler.RespondNotFound(c, "04", "not found", "image not found", "")
		return
	}

	var dir = fmt.Sprintf("%s/src/gitlab.com/bri-bridge/backend-bridge-api/storage/email/logo.png",
		os.Getenv("GOPATH"))
	c.File(dir)
	return
}
