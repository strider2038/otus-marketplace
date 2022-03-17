package di

import (
	"notification-service/internal/kafka"
	"notification-service/internal/messaging"
	"notification-service/internal/postgres"
	"notification-service/internal/postgres/database"

	"github.com/jackc/pgx/v4/pgxpool"
	segmentio "github.com/segmentio/kafka-go"
)

func NewIdentityConsumer(connection *pgxpool.Pool, config Config) *kafka.Consumer {
	db := database.New(connection)
	users := postgres.NewUserRepository(db)

	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers: []string{config.KafkaConsumerURL},
		GroupID: "notification",
		Topic:   "identity-events",
	})

	consumer := kafka.NewConsumer(reader, map[string]kafka.Processor{
		"Identity/UserCreated": messaging.NewUserCreatedProcessor(users),
	})

	return consumer
}
