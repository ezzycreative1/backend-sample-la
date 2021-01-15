package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	routes "backend-sample-la/routes"
)

func main() {
	var err error
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
	}))

	routes.HealthCheckHTTPHandler(router)

	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	// Server
	if err := router.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))); err != nil {
		log.Fatal(err)
	}
}
