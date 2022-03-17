package di

import (
	"trading-service/internal/kafka"
	"trading-service/internal/messaging"
	"trading-service/internal/postgres"
	"trading-service/internal/trading"

	"github.com/muonsoft/validation"
	segmentio "github.com/segmentio/kafka-go"
	"github.com/strider2038/pkg/persistence/pgx"
)

type Container struct {
	connection pgx.Connection
	config     Config

	billingDispatcher       *kafka.Dispatcher
	tradingDispatcher       *kafka.Dispatcher
	purchaseOrderRepository *postgres.PurchaseOrderRepository
	sellOrderRepository     *postgres.SellOrderRepository
	itemRepository          *postgres.ItemRepository
	txManager               *pgx.TransactionManager
	billingAdapter          *messaging.BillingAdapter
	tradingAdapter          *messaging.TradingAdapter
	dealer                  *trading.Dealer
	validator               *validation.Validator
}

func NewContainer(connection pgx.Connection, config Config) (*Container, error) {
	var err error

	c := &Container{connection: connection, config: config}
	c.billingDispatcher = kafka.NewDispatcher(&segmentio.Writer{
		Addr:     segmentio.TCP(config.KafkaProducerURL),
		Topic:    "billing-commands",
		Balancer: &segmentio.LeastBytes{},
	})
	c.tradingDispatcher = kafka.NewDispatcher(&segmentio.Writer{
		Addr:     segmentio.TCP(config.KafkaProducerURL),
		Topic:    "trading-events",
		Balancer: &segmentio.LeastBytes{},
	})

	c.purchaseOrderRepository = postgres.NewPurchaseOrderRepository(connection)
	c.sellOrderRepository = postgres.NewSellOrderRepository(connection)
	c.itemRepository = postgres.NewItemRepository(connection)
	c.txManager = pgx.NewTransactionManager(connection)
	c.billingAdapter = messaging.NewBillingAdapter(c.billingDispatcher)
	c.tradingAdapter = messaging.NewTradingAdapter(c.tradingDispatcher)
	c.dealer = trading.NewDealer(
		c.itemRepository,
		c.purchaseOrderRepository,
		c.sellOrderRepository,
		c.txManager,
		c.billingAdapter,
		c.tradingAdapter,
	)
	c.validator, err = validation.NewValidator()
	if err != nil {
		return nil, err
	}

	return c, nil
}
