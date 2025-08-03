package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	pb "oprosdom.ru/shared/models/proto"
)

type Producer struct {
	writer *kafka.Writer
}

func NewKafka(ctx context.Context, conn string, topic string) (*Producer, error) {

	// проверяем соединение перед инициализацией
	if _, err := kafka.DialContext(ctx, "tcp", conn); err != nil {
		return nil, fmt.Errorf("kafka connection failed: %w", err)
	}

	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(conn),
			Topic:        topic,
			Balancer:     &kafka.Hash{},         // обязательно Hash чтоб работало Партиционирование по номеру телефона
			BatchTimeout: 10 * time.Millisecond, // это сколько времени ждать пока накопятся сообщения для отправки чтоб отправить оптом
		},
	}, nil
}

func (p *Producer) Send(ctx context.Context, msg *pb.MsgCode) error {
	protoData, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf msg: %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctxTimeout, kafka.Message{
		Key:   []byte(msg.PhoneNumber), // Партиционирование по номеру телефона упрощает обработку и масштабирование, сохраняя последовательность для каждого пользователя. Если юзер несколько кодов запросит, чтоб и он и мы понимали где последний актуальный код
		Value: protoData,
	})

	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
