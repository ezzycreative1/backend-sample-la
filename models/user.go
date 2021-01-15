package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// Users ..
type Users struct {
	User []User `json:"users"`
}

// User is struct or model data user
type User struct {
	Id          string         `gorm:"primary_key;type:uuid";"AUTO_INCREMENT" json:"id"`
	Name        string         `gorm:"not null" json:"name"`         //`gorm:"type:varchar(200)”`
	Email       string         `gorm:"not null,unique" json:"email"` //`gorm:"type:varchar(100)”;"unique"`
	Password    string         `gorm:"not null" json:",omitempty"`
	PhoneNumber string         `json:"phone_number"`
	RoleId      int64          `json:"role_id,omitempty"` //`gorm:"ForeignKey:id,not null"` //`gorm:"type:varchar(20)”`
	Photo       string         `json:"photo"`
	VerifyToken string         `json:",omitempty"`
	VerifedAt   string         `json:",omitempty"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// Role is struct or model data role
type Role struct {
	Id        int64          `gorm:"primary_key";"AUTO_INCREMENT"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// UserResponse ...
type UserResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// Profile ..
type Profile struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hd        string `json:"hd"`
	Locale    string `json:"locale"`
	Name      string `json:"name"`
	Picture   struct {
		Data struct {
			Url string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

// AccessToken ..
type AccessToken struct {
	Token  string
	Expiry int64
	Name   string
	Id     string
}

// this for response

type ListUser struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RoleId      int64  `json:"role_id"`
}
