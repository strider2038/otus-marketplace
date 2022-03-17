package di

import (
	"trading-service/internal/kafka"
	"trading-service/internal/messaging"

	segmentio "github.com/segmentio/kafka-go"
)

func NewBillingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers: []string{c.config.KafkaConsumerURL},
		GroupID: "trading",
		Topic:   "billing-commands",
	})

	return kafka.NewConsumer(reader, map[string]kafka.Processor{
		"Billing/PaymentSucceeded": messaging.NewPaymentSucceededProcessor(c.dealer),
		"Billing/PaymentDeclined":  messaging.NewPaymentDeclinedProcessor(c.dealer),
		"Billing/AccrualApproved":  messaging.NewAccrualApprovedProcessor(c.dealer),
	})
}
