package kafka

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	// svc    service.MessageService
}

// func NewConsumer(brokers []string, topic, groupID string, svc service.MessageService) *Consumer {
func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,
			MaxBytes:       1e6,
			CommitInterval: time.Second,
			MaxWait:        5 * time.Second,
		}),
		// svc: svc,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	defer c.reader.Close()

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "Stopping Kafka consumer")
			return nil
		default:
			// Читаем сообщение с таймаутом
			//msg, err := c.reader.ReadMessage(ctx)
			_, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				slog.ErrorContext(ctx, "Error reading message", "error", err)
				continue
			}

			log.Println("got message!")

			// // Декодируем сообщение
			// var message service.Message
			// if err := json.Unmarshal(msg.Value, &message); err != nil {
			// 	slog.ErrorContext(ctx, "Error decoding message",
			// 		"offset", msg.Offset, "partition", msg.Partition, "error", err)
			// 	continue
			// }

			// // Обрабатываем сообщение через сервис
			// if err := c.svc.ProcessMessage(ctx, message); err != nil {
			// 	slog.ErrorContext(ctx, "Error processing message",
			// 		"phone", message.Phone, "error", err)
			// }
		}
	}
}
