package di

import (
	"context"
	"log"
	"os"

	"trading-service/internal/kafka"
	"trading-service/internal/messaging"
	"trading-service/internal/postgres"
	redisadapter "trading-service/internal/redis"
	"trading-service/internal/trading"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
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
	config Config

	dbConnection    pgxadapter.Connection
	redisConnection *redis.Client
	logger          zerolog.Logger

	locker                  *redisadapter.LockerAdapter
	billingDispatcher       *kafka.Dispatcher
	tradingDispatcher       *kafka.Dispatcher
	purchaseOrderRepository *postgres.PurchaseOrderRepository
	sellOrderRepository     *postgres.SellOrderRepository
	itemRepository          *postgres.ItemRepository
	userItemRepository      *postgres.UserItemRepository
	txManager               *pgxadapter.TransactionManager
	billingAdapter          *messaging.BillingAdapter
	tradingAdapter          *messaging.TradingAdapter
	dealer                  *trading.Dealer
	validator               *validation.Validator
	purchaseState           *redisadapter.StateRepository
	sellState               *redisadapter.StateRepository
}

func NewContainer(config Config) (*Container, error) {
	var err error

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
	c.logger = zerolog.New(os.Stdout)
	c.locker = redisadapter.NewLockerAdapter(redislock.New(redisClient))

	c.billingDispatcher = kafka.NewDispatcher(
		&segmentio.Writer{
			Addr:        segmentio.TCP(config.KafkaProducerURL),
			Topic:       "billing-commands",
			Balancer:    &segmentio.LeastBytes{},
			ErrorLogger: log.Default(),
		},
		c.logger,
	)
	c.tradingDispatcher = kafka.NewDispatcher(
		&segmentio.Writer{
			Addr:        segmentio.TCP(config.KafkaProducerURL),
			Topic:       "trading-events",
			Balancer:    &segmentio.LeastBytes{},
			ErrorLogger: log.Default(),
		},
		c.logger,
	)

	c.purchaseOrderRepository = postgres.NewPurchaseOrderRepository(c.dbConnection)
	c.sellOrderRepository = postgres.NewSellOrderRepository(c.dbConnection)
	c.itemRepository = postgres.NewItemRepository(c.dbConnection)
	c.userItemRepository = postgres.NewUserItemRepository(c.dbConnection)
	c.txManager = pgxadapter.NewTransactionManager(c.dbConnection)
	c.billingAdapter = messaging.NewBillingAdapter(c.billingDispatcher)
	c.tradingAdapter = messaging.NewTradingAdapter(c.tradingDispatcher)
	c.dealer = trading.NewDealer(
		c.itemRepository,
		c.userItemRepository,
		c.purchaseOrderRepository,
		c.sellOrderRepository,
		c.txManager,
		c.billingAdapter,
		c.tradingAdapter,
	)

	c.sellState = redisadapter.NewStateRepository("sell", c.redisConnection, c.config.StateTimeout)
	c.purchaseState = redisadapter.NewStateRepository("purchase", c.redisConnection, c.config.StateTimeout)

	c.validator, err = validation.NewValidator()
	if err != nil {
		return nil, err
	}

	return c, nil
}
