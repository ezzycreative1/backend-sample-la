package migration

import (
	"fmt"
	"log"

	"backend-sample-la/database"
	"backend-sample-la/models"
)

// RunRollback ..
func RunRollback() {

	db := database.PostsqlConn()

	if db.Error != nil {
		log.Fatalln(db.Error.Error())
	}

	if exist := db.HasTable("users"); exist {
		fmt.Println("drop table users")
		err := db.DropTable("users")
		if err == nil {
			fmt.Println("success drop table users")
		}
	}

	if exist := db.HasTable("roles"); exist {
		fmt.Println("drop table roles")
		err := db.DropTable("roles")
		if err == nil {
			fmt.Println("success drop table roles")
		}
	}
}

// RunMigration ..
func RunMigration() {
	db := database.PostsqlConn()

	if db.Error != nil {
		log.Fatalln(db.Error.Error())
	}

	if exist := db.HasTable("roles"); !exist {
		fmt.Println("migrate table roles")
		err := db.CreateTable(&models.Role{})
		if err == nil {
			fmt.Println("success migrate table roles")
		}
	}

	if exist := db.HasTable("users"); !exist {
		fmt.Println("migrate table users")
		err := db.CreateTable(&models.User{})
		if err == nil {
			fmt.Println("success migrate table users")
		}
	}

}
