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

func NewTradingConsumer(connection *pgxpool.Pool, config Config) *kafka.Consumer {
	db := database.New(connection)
	users := postgres.NewUserRepository(db)
	notifications := postgres.NewNotificationRepository(db)

	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{config.KafkaConsumerURL},
		GroupID:     "notification",
		Topic:       "trading-events",
		ErrorLogger: log.Default(),
	})

	return kafka.NewConsumer(reader, zerolog.New(os.Stdout), map[string]kafka.Processor{
		messaging.DealSucceeded{}.Name():  messaging.NewDealSucceededProcessor(users, notifications),
		messaging.PurchaseFailed{}.Name(): messaging.NewPurchaseFailedProcessor(users, notifications),
	})
}
