package di

import (
	"log"

	"history-service/internal/kafka"
	"history-service/internal/messaging"
	"history-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewTradingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "history",
		Topic:       "trading-events",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		messaging.DealSucceeded{}.Name(): messaging.NewDealSucceededProcessor(c.dealRepository),
	})

	return kafka.NewConsumer(reader, c.logger, monitoring.NewMessagingMiddleware(mux, c.metrics))
}
