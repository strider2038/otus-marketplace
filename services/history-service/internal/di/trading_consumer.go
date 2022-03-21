package di

import (
	"log"
	"os"

	"history-service/internal/kafka"
	"history-service/internal/messaging"

	"github.com/rs/zerolog"
	segmentio "github.com/segmentio/kafka-go"
)

func NewTradingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "history",
		Topic:       "trading-events",
		ErrorLogger: log.Default(),
	})

	return kafka.NewConsumer(reader, zerolog.New(os.Stdout), map[string]kafka.Processor{
		messaging.DealSucceeded{}.Name(): messaging.NewDealSucceededProcessor(c.dealRepository),
	})
}
