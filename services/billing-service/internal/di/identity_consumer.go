package di

import (
	"billing-service/internal/kafka"
	"billing-service/internal/messaging"
	"billing-service/internal/postgres"

	segmentio "github.com/segmentio/kafka-go"
	"github.com/strider2038/pkg/persistence/pgx"
)

func NewIdentityConsumer(connection pgx.Connection, config Config) *kafka.Consumer {
	accounts := postgres.NewAccountRepository(connection)

	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers: []string{config.KafkaConsumerURL},
		GroupID: "billing",
		Topic:   "identity-events",
	})

	consumer := kafka.NewConsumer(reader, map[string]kafka.Processor{
		"Identity/UserCreated": messaging.NewUserCreatedProcessor(accounts),
	})

	return consumer
}
