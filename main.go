package main

import (
	"fmt"
	"log"
	"os"

	"backend-sample-la/database/seeder"

	UserRepo "backend-sample-la/app/user/repository"
	UserUsecase "backend-sample-la/app/user/usecase"
	"backend-sample-la/database"
	"backend-sample-la/database/migration"
	middle "backend-sample-la/middleware"
	routes "backend-sample-la/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	err error
)

func main() {

	err = godotenv.Load()

	db := database.PostsqlConn()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	dbEvent := os.Getenv("DBEVENT")

	if dbEvent == "rollback_migrate" || dbEvent == "rollback" {
		migration.RunRollback()
	}

	if dbEvent == "migrate" || dbEvent == "rollback_migrate" {
		migration.RunMigration()
		seeder.RoleSeeder()
		seeder.UserSeeder()
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Init Usecase and Repository
	// User
	userRepo := UserRepo.NewUserRepository(db)
	userUsecase := UserUsecase.NewUserUsecase(userRepo)

	routes.NewHandlerWithoutCORS(router, userUsecase)

	// CORS
	router.Use(CORSMiddleware())

	// Size Images
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Middleware Token
	router.Use(middle.AuthMiddleware())

	// Health check
	routes.HealthCheckHTTPHandler(router)

	// User
	routes.NewUserHandler(router, userUsecase)

	// Server
	if err := router.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))); err != nil {
		log.Fatal(err)
	}
}

// CORSMiddleware ..
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, secret-key, BRIDGE-Nonce, BRIDGE-Signature")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			//c.Next()
			return
		}

		c.Next()
	}
}
