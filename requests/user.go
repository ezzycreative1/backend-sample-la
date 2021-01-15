package requests

// RegisterRequest ..
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest ..
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser ..
type CreateUser struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	NoHp   string `json:"noHp"`
	RoleId int    `json:"roleId"`
}

// UpdateUser ..
type UpdateUser struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	NoHp   string `json:"noHp"`
	RoleId int    `json:"roleId"`
}

// ChangePhoto ..
type ChangePhoto struct {
	Id    string `json:"id"`
	Photo string `json:"photo"`
}

// ChangeName ..
type ChangeName struct {
	Id   string `json:"Id"`
	Name string `json:"name"`
}
type ForgotPassword struct {
	Email string `json:"email"`
}

type ValidateGoogleToken struct {
	Token string `json:"token" validate:"required"`
}

// ChangePassword ..
type ChangePassword struct {
	UserId          string
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
