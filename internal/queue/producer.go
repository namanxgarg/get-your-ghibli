package queue

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"os"
)

var ch *amqp091.Channel

func InitQueue() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic("❌ Failed to connect to RabbitMQ: " + err.Error())
	}
	channel, err := conn.Channel()
	if err != nil {
		panic("❌ Failed to open channel: " + err.Error())
	}

	err = channel.QueueDeclare(
		"ghibli_tasks", false, false, false, false, nil,
	)
	if err != nil {
		panic("❌ Failed to declare queue: " + err.Error())
	}
	ch = channel
	fmt.Println("✅ Connected to RabbitMQ")
}

type Task struct {
	UploadID uint   `json:"upload_id"`
	PhotoURL string `json:"photo_url"`
}

func SendTask(uploadID uint, photoURL string) error {
	task := Task{UploadID: uploadID, PhotoURL: photoURL}
	body, _ := json.Marshal(task)

	err := ch.Publish(
		"", "ghibli_tasks", false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}
