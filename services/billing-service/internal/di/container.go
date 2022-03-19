package di

import (
	"log"
	"os"

	"billing-service/internal/billing"
	"billing-service/internal/kafka"
	"billing-service/internal/postgres"

	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	segmentio "github.com/segmentio/kafka-go"
	"github.com/strider2038/pkg/persistence/pgx"
)

type Container struct {
	connection pgx.Connection
	config     Config
	logger     zerolog.Logger

	accountRepository   *postgres.AccountRepository
	operationRepository *postgres.OperationRepository
	txManager           *pgx.TransactionManager
	billingDispatcher   *kafka.Dispatcher
	broker              *billing.BrokerAccount
	validator           *validation.Validator
}

func NewContainer(connection pgx.Connection, config Config) (*Container, error) {
	brokerID, err := uuid.FromString(config.BrokerID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid broker id")
	}

	c := &Container{connection: connection, config: config}

	c.logger = zerolog.New(os.Stdout)
	c.accountRepository = postgres.NewAccountRepository(connection)
	c.operationRepository = postgres.NewOperationRepository(connection)
	c.txManager = pgx.NewTransactionManager(connection)
	c.broker = billing.NewBrokerAccount(
		brokerID,
		c.accountRepository,
		c.operationRepository,
	)
	c.billingDispatcher = kafka.NewDispatcher(
		&segmentio.Writer{
			Addr:        segmentio.TCP(config.KafkaProducerURL),
			Topic:       "billing-events",
			Balancer:    &segmentio.LeastBytes{},
			ErrorLogger: log.Default(),
		},
		c.logger,
	)
	c.validator, err = validation.NewValidator()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create validator")
	}

	return c, nil
}
