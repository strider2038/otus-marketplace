package di

import (
	"log"
	"os"

	"statistics-service/internal/kafka"
	"statistics-service/internal/messaging"

	"github.com/rs/zerolog"
	segmentio "github.com/segmentio/kafka-go"
)

func NewTradingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "statistics",
		Topic:       "trading-events",
		ErrorLogger: log.Default(),
	})

	return kafka.NewConsumer(reader, zerolog.New(os.Stdout), map[string]kafka.Processor{
		messaging.DealSucceeded{}.Name(): messaging.NewDealSucceededProcessor(
			c.dailyDeals,
			c.totalDailyDeals,
			c.top10Deals,
		),
	})
}
