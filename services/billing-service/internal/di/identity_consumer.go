package di

import (
	"log"

	"billing-service/internal/kafka"
	"billing-service/internal/messaging"
	"billing-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewIdentityConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "billing",
		Topic:       "identity-events",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		"Identity/UserCreated": messaging.NewUserCreatedProcessor(c.accountRepository),
	})

	return kafka.NewConsumer(
		reader,
		c.logger,
		monitoring.NewMessagingMiddleware(mux, c.metrics),
	)
}
