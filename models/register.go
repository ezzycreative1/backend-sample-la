package models

import "time"

// RegisterResponse ...
type RegisterResponse struct {
	Message string `json:"message"`
}

// RegisterDocument ..
type RegisterDocument struct {
	ID          string `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        int    `json:"role"`
	VerifyToken string
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}
