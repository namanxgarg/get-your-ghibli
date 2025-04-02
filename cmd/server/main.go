package main

import (
	"log"
	"os"

	"get-your-ghibli/internal/auth"
	"get-your-ghibli/internal/db"
	"get-your-ghibli/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, falling back to system env")
	}

	db.Init()
	auth.InitRedis(os.Getenv("REDIS_URL"))
	db.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	// Simple healthcheck route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/auth/request-otp", auth.RequestOTPHandler)

	log.Println("🚀 Server is running at http://localhost:8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
