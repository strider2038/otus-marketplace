package di

import (
	"log"

	"trading-service/internal/kafka"
	"trading-service/internal/messaging"
	"trading-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewBillingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "trading",
		Topic:       "billing-events",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		"Billing/PaymentSucceeded": messaging.NewPaymentSucceededProcessor(c.dealer),
		"Billing/PaymentDeclined":  messaging.NewPaymentDeclinedProcessor(c.dealer),
		"Billing/AccrualApproved":  messaging.NewAccrualApprovedProcessor(c.dealer),
	})

	return kafka.NewConsumer(
		reader,
		c.logger,
		monitoring.NewMessagingMiddleware(mux, c.metrics),
	)
}
