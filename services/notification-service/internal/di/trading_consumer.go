package di

import (
	"log"

	"notification-service/internal/kafka"
	"notification-service/internal/messaging"
	"notification-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewTradingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "notification",
		Topic:       "trading-events",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		messaging.DealSucceeded{}.Name():  messaging.NewDealSucceededProcessor(c.users, c.notifications),
		messaging.PurchaseFailed{}.Name(): messaging.NewPurchaseFailedProcessor(c.users, c.notifications),
	})

	return kafka.NewConsumer(reader, c.logger, monitoring.NewMessagingMiddleware(mux, c.metrics))
}
