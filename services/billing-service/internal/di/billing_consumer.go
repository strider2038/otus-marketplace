package di

import (
	"billing-service/internal/kafka"
	"billing-service/internal/messaging"
	"billing-service/internal/postgres"

	segmentio "github.com/segmentio/kafka-go"
	"github.com/strider2038/pkg/persistence/pgx"
)

func NewBillingConsumer(connection pgx.Connection, config Config) *kafka.Consumer {
	accounts := postgres.NewAccountRepository(connection)
	payments := postgres.NewOperationRepository(connection)
	txManager := pgx.NewTransactionManager(connection)

	writer := &segmentio.Writer{
		Addr:     segmentio.TCP(config.KafkaProducerURL),
		Topic:    "billing-events",
		Balancer: &segmentio.LeastBytes{},
	}
	dispatcher := kafka.NewDispatcher(writer)

	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers: []string{config.KafkaConsumerURL},
		GroupID: "billing",
		Topic:   "billing-commands",
	})

	consumer := kafka.NewConsumer(reader, map[string]kafka.Processor{
		"Billing/CreatePayment": messaging.NewCreatePaymentProcessor(accounts, payments, txManager, dispatcher),
		"Billing/CreateAccrual": messaging.NewCreateAccrualProcessor(accounts, payments, txManager, dispatcher),
	})

	return consumer
}
