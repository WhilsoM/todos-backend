package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type UserRegistered struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "group-test",
		GroupID: "analytic-service",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		var user UserRegistered

		if err := json.Unmarshal(msg.Value, &user); err != nil {
			log.Fatalf("failed to unmarshal value from consumer %v %v %v", msg.Value, reader.Config().Topic, reader.Config().GroupID)
		}

		log.Println(user.Email, user.UserID)

		log.Printf(
			"message=%s partition=%d offset=%d",
			string(msg.Value),
			msg.Partition,
			msg.Offset,
		)
	}
}
