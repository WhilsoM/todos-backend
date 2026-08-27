package events

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type UserEventProducer struct {
	writer *kafka.Writer
}

func NewUserEventProducer(broker string) *UserEventProducer {
	return &UserEventProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        "user-events",
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (u *UserEventProducer) PublishUserRegistered(ctx context.Context, userID int) error {
	return u.writer.WriteMessages(ctx, kafka.Message{
		Value: fmt.Appendf([]byte("user registered"), "id %d", userID),
	})
}

func (u *UserEventProducer) Close() error {
	return u.writer.Close()
}
