package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var (
	psqlConn *gorm.DB
	err      error
)

// initialize database
func init() {
	setupPostsqlConn()
}

func setupPostsqlConn() {

	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	dns := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("POSTGRE_HOST"), os.Getenv("POSTGRE_PORT"), os.Getenv("POSTGRE_USER"), os.Getenv("POSTGRE_DBNAME"), os.Getenv("POSTGRE_PASS"))
	psqlConn, err = gorm.Open("postgres", dns)

	err = psqlConn.DB().Ping()
	if err != nil {
		panic(err)
	}

	psqlConn.LogMode(true)

}

// PostsqlConn return mysql connection from gorm ORM
func PostsqlConn() *gorm.DB {
	return psqlConn
}
