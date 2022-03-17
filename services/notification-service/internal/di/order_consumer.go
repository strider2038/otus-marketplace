package di

import (
	"notification-service/internal/kafka"
	"notification-service/internal/messaging"
	"notification-service/internal/postgres"
	"notification-service/internal/postgres/database"

	"github.com/jackc/pgx/v4/pgxpool"
	segmentio "github.com/segmentio/kafka-go"
)

func NewOrderConsumer(connection *pgxpool.Pool, config Config) *kafka.Consumer {
	db := database.New(connection)
	users := postgres.NewUserRepository(db)
	notifications := postgres.NewNotificationRepository(db)

	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers: []string{config.KafkaConsumerURL},
		GroupID: "notification",
		Topic:   "order-events",
	})

	consumer := kafka.NewConsumer(reader, map[string]kafka.Processor{
		"Ordering/OrderSucceeded": messaging.NewOrderSucceededProcessor(users, notifications),
		"Ordering/OrderFailed":    messaging.NewOrderFailedProcessor(users, notifications),
	})

	return consumer
}
