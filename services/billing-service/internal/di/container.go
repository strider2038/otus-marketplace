package di

import (
	"context"
	"log"
	"os"

	"billing-service/internal/billing"
	"billing-service/internal/kafka"
	"billing-service/internal/postgres"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muonsoft/validation"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	segmentio "github.com/segmentio/kafka-go"
	pgxadapter "github.com/strider2038/pkg/persistence/pgx"
)

type Container struct {
	dbConnection    pgxadapter.Connection
	redisConnection *redis.Client
	config          Config
	logger          zerolog.Logger

	locker              *redislock.Client
	accountRepository   *postgres.AccountRepository
	operationRepository *postgres.OperationRepository
	txManager           *pgxadapter.TransactionManager
	billingDispatcher   *kafka.Dispatcher
	broker              *billing.BrokerAccount
	validator           *validation.Validator
}

func NewContainer(config Config) (*Container, error) {
	brokerID, err := uuid.FromString(config.BrokerID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid broker id")
	}

	dbConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse postgres config:")
	}

	logger := zerolog.New(os.Stdout)
	dbConfig.ConnConfig.Logger = zerologadapter.NewLogger(logger)
	dbConfig.ConnConfig.LogLevel = pgx.LogLevelInfo

	pool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	redisClient := redis.NewClient(&redis.Options{Network: "tcp", Addr: config.RedisURL})
	ping := redisClient.Ping(context.Background())
	err = ping.Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to redis")
	}

	c := &Container{config: config}

	c.dbConnection = pgxadapter.NewPool(pool)
	c.redisConnection = redisClient
	c.locker = redislock.New(c.redisConnection)
	c.logger = zerolog.New(os.Stdout)
	c.accountRepository = postgres.NewAccountRepository(c.dbConnection)
	c.operationRepository = postgres.NewOperationRepository(c.dbConnection)
	c.txManager = pgxadapter.NewTransactionManager(c.dbConnection)
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
