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
	user := UserRegistered{
		UserID: 1,
		Email:  "test@gmail.com",
	}
	ctx := context.Background()
	bytes, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("failed to marshal user struct %v, %v", err, bytes)
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "group-test",
	}
	defer writer.Close()

	kafkaMsg := kafka.Message{
		Value: bytes,
	}
	if err := writer.WriteMessages(ctx,
		kafkaMsg,
	); err != nil {
		log.Fatal(err)

	}

	log.Println("message sent")
}
