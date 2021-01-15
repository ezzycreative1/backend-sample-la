package seeder

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gitlab.com/bri-bridge/backend-bridge-api/database"
	"gitlab.com/bri-bridge/backend-bridge-api/models"
	"golang.org/x/crypto/bcrypt"
)

func UserSeeder() {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatalln(err.Error())
	}
	now := time.Now()
	user := models.User{
		Id:        u.String(),
		Name:      "Super User",
		Email:     "superuser@bridge.com",
		Password:  hashPassword("Superuser123"),
		RoleId:    1,
		VerifedAt: fmt.Sprintf("%v", now),
		CreatedAt: now,
		UpdatedAt: now,
	}

	db := database.PostsqlConn()
	if err := db.Create(&user).Error; err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("superuser has been created\n")
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Fatalln(err.Error())
	}
	return string(bytes)
}
