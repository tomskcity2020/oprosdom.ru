package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"oprosdom.ru/msvc_codesender/internal/models"
	"oprosdom.ru/msvc_codesender/internal/service"
	"oprosdom.ru/shared/models/pb"
)

type Consumer struct {
	reader *kafka.Reader
	svc    service.ServiceInterface
}

func NewConsumer(brokers []string, topic string, groupID string, svc service.ServiceInterface) *Consumer {
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
		svc: svc,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	defer func() {
		log.Println("Kafka reader closed")
		if err := c.reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Println("Kafka consumer stopped by context cancel")
				return nil
			}

			log.Printf("Error reading message: %v", err)
			continue
		}

		// Декодируем protobuf сообщение
		var msgCode pb.MsgCode
		if err := proto.Unmarshal(msg.Value, &msgCode); err != nil {
			log.Printf("Error decoding protobuf message [offset:%d partition:%d]: %v", msg.Offset, msg.Partition, err)
			continue
		}

		// Преобразуем в структуру для дальнейшей первичной обработки
		unsafeMsg := models.UnsafeMsg{
			Phone: msgCode.GetPhone(),
			Code:        msgCode.GetCode(),
			Retry:       msgCode.GetRetry(),
		}

		// Первичную проверку нужно проводить на уровне хендлера и отдаем в сервис уже valid, здесь аналогично
		validMsg, err := unsafeMsg.Validate()
		if err != nil {
			log.Println("Error validating msg")
		}

		// Обрабатываем сообщение
		if err := c.svc.AddMessage(ctx, validMsg); err != nil {
			log.Printf("Error processing message: %v", err)
		}

	}
}
