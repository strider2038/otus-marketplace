package di

import (
	"log"

	"billing-service/internal/kafka"
	"billing-service/internal/messaging"
	"billing-service/internal/monitoring"

	segmentio "github.com/segmentio/kafka-go"
)

func NewBillingConsumer(c *Container) *kafka.Consumer {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:     []string{c.config.KafkaConsumerURL},
		GroupID:     "billing",
		Topic:       "billing-commands",
		ErrorLogger: log.Default(),
	})

	mux := kafka.NewMux(c.logger, map[string]kafka.Processor{
		"Billing/CreatePayment": messaging.NewCreatePaymentProcessor(
			c.accountRepository,
			c.operationRepository,
			c.broker,
			c.txManager,
			c.billingDispatcher,
		),
		"Billing/CreateAccrual": messaging.NewCreateAccrualProcessor(
			c.accountRepository,
			c.operationRepository,
			c.broker,
			c.txManager,
			c.billingDispatcher,
		),
	})

	return kafka.NewConsumer(
		reader,
		c.logger,
		monitoring.NewMessagingMiddleware(mux, c.metrics),
	)
}
