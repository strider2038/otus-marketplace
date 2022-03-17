package di

import (
	"billing-service/internal/kafka"
	"billing-service/internal/messaging"

	segmentio "github.com/segmentio/kafka-go"
)

func NewIdentityConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers: []string{c.config.KafkaConsumerURL},
		GroupID: "billing",
		Topic:   "identity-events",
	})

	return kafka.NewConsumer(reader, map[string]kafka.Processor{
		"Identity/UserCreated": messaging.NewUserCreatedProcessor(c.accountRepository),
	})
}
