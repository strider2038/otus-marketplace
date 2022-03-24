package di

import (
	"log"

	"statistics-service/internal/kafka"
	"statistics-service/internal/messaging"
	"statistics-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewTradingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "statistics",
		Topic:       "trading-events",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		messaging.DealSucceeded{}.Name(): messaging.NewDealSucceededProcessor(
			c.dailyDeals,
			c.totalDailyDeals,
			c.top10Deals,
		),
	})

	return kafka.NewConsumer(reader, c.logger, monitoring.NewMessagingMiddleware(mux, c.metrics))
}
