package di

import (
	"log"

	"notification-service/internal/kafka"
	"notification-service/internal/messaging"
	"notification-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewIdentityConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "notification",
		Topic:       "identity-events",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		messaging.UserCreated{}.Name(): messaging.NewUserCreatedProcessor(c.users),
	})

	return kafka.NewConsumer(reader, c.logger, monitoring.NewMessagingMiddleware(mux, c.metrics))
}
