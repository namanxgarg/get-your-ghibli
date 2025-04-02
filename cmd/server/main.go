package main

import (
	"log"
	"os"

	"get-your-ghibli/internal/auth"
	"get-your-ghibli/internal/db"
	"get-your-ghibli/internal/upload"
	"get-your-ghibli/pkg/models"

	"get-your-ghibli/internal/queue"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"get-your-ghibli/internal/gallery"
	"get-your-ghibli/internal/order"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, falling back to system env")
	}

	db.Init()
	auth.InitRedis(os.Getenv("REDIS_URL"))
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Upload{})

	db.DB.AutoMigrate(&models.GeneratedImage{})
	db.DB.AutoMigrate(&models.Order{})

	r := gin.Default()

	// Simple healthcheck route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/auth/request-otp", auth.RequestOTPHandler)
	r.POST("/auth/verify-otp", auth.VerifyOTPHandler)

	r.POST("/webhook/razorpay", upload.RazorpayWebhookHandler)

	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	protected.GET("/auth/me", auth.MeHandler)
	protected.GET("/gallery", gallery.GalleryHandler)


	protected.POST("/upload", upload.UploadHandler)
	protected.POST("/order", order.PlaceOrderHandler)


	queue.InitQueue()

	log.Println("🚀 Server is running at http://localhost:8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
