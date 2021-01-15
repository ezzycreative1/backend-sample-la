package seeder

import (
	"fmt"
	"log"
	"time"

	"backend-sample-la/database"
	"backend-sample-la/models"
)

func RoleSeeder() {

	db := database.PostsqlConn()

	now := time.Now()

	var roles []models.Role
	// superuser
	var superuser = models.Role{
		Id:        1,
		Name:      "superuser",
		CreatedAt: now,
		UpdatedAt: now,
	}
	roles = append(roles, superuser)

	var admin = models.Role{
		Id:        2,
		Name:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}
	roles = append(roles, admin)

	var user = models.Role{
		Id:        3,
		Name:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	roles = append(roles, user)

	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Printf("role %s has been created\n", role.Name)
	}

}
