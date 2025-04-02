package main

import (
	"encoding/json"
	"fmt"
	"get-your-ghibli/internal/db"
	"get-your-ghibli/pkg/models"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

type Task struct {
	UploadID uint   `json:"upload_id"`
	PhotoURL string `json:"photo_url"`
}

func main() {
	godotenv.Load()
	db.Init()

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic("❌ RabbitMQ connection failed: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("❌ Failed to open channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"ghibli_tasks",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		panic("❌ Failed to declare queue: " + err.Error())
	}
	

	msgs, err := ch.Consume(
		"ghibli_tasks", "", true, false, false, false, nil,
	)
	if err != nil {
		panic("❌ Failed to consume queue: " + err.Error())
	}

	fmt.Println("🎨 Ghibli Generator Worker started...")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var task Task
			json.Unmarshal(d.Body, &task)
			fmt.Printf("🎯 Received task for upload #%d\n", task.UploadID)

			// Simulate AI generation
			generateGhibliImages(task.UploadID)
		}
	}()
	<-forever
}

func generateGhibliImages(uploadID uint) {
	// Simulate delay
	time.Sleep(5 * time.Second)

	// Create 10 fake URLs
	for i := 1; i <= 10; i++ {
		img := models.GeneratedImage{
			UploadID: uploadID,
			ImageURL: fmt.Sprintf("https://fake-ghibli.s3.com/%d_%d.png", uploadID, rand.Intn(10000)),
		}
		db.DB.Create(&img)
	}

	db.DB.Model(&models.Upload{}).Where("id = ?", uploadID).Update("status", "ready")
	fmt.Printf("✅ Upload #%d marked as READY with 10 Ghibli images!\n", uploadID)
}
