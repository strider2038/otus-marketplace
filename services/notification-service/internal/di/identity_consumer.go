package di

import (
	"log"
	"os"

	"notification-service/internal/kafka"
	"notification-service/internal/messaging"
	"notification-service/internal/postgres"
	"notification-service/internal/postgres/database"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	segmentio "github.com/segmentio/kafka-go"
)

func NewIdentityConsumer(connection *pgxpool.Pool, config Config) *kafka.Consumer {
	db := database.New(connection)
	users := postgres.NewUserRepository(db)

	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{config.KafkaConsumerURL},
		GroupID:     "notification",
		Topic:       "identity-events",
		ErrorLogger: log.Default(),
	})

	return kafka.NewConsumer(reader, zerolog.New(os.Stdout), map[string]kafka.Processor{
		messaging.UserCreated{}.Name(): messaging.NewUserCreatedProcessor(users),
	})
}
