package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-events",
		GroupID: "user-events",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"message=%s partition=%d offset=%d",
			string(msg.Value),
			msg.Partition,
			msg.Offset,
		)
	}
}
