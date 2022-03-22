package apitest

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"

	"trading-service/internal/api"
	"trading-service/internal/messaging"
	"trading-service/internal/trading"
	"trading-service/test/mock"

	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/stretchr/testify/suite"
	"github.com/strider2038/pkg/persistence"
)

var (
	userID      = uuid.FromStringOrNil("970a9f0b-ad00-4f65-8d93-aaa00742b0ef")
	itemID      = uuid.FromStringOrNil("31e2951a-08da-4ae6-b37b-928cf2507448")
	userItemID  = uuid.FromStringOrNil("cf6f9f52-4290-45a9-9773-d3af268bedf4")
	brokerID    = uuid.FromStringOrNil("432f843a-446a-49ad-8764-eca9d2ac64b8")
	sellerID    = uuid.FromStringOrNil("771aea75-b73c-4fb6-b9d4-5a6a018a0d04")
	purchaserID = uuid.FromStringOrNil("e9581c64-74c2-492f-85cd-04f5abeeabe8")
	sellOrderID = uuid.FromStringOrNil("227ec860-64d4-477c-8fd6-b482b6dca06e")
)

type APISuite struct {
	suite.Suite

	items          *mock.ItemRepository
	userItems      *mock.UserItemRepository
	purchaseOrders *mock.PurchaseOrderRepository
	sellOrders     *mock.SellOrderRepository
	dispatcher     *mock.MessageDispatcher

	api http.Handler
}

func (suite *APISuite) SetupTest() {
	suite.items = mock.NewItemRepository()
	suite.userItems = mock.NewUserItemRepository()
	suite.purchaseOrders = mock.NewPurchaseOrderRepository()
	suite.sellOrders = mock.NewSellOrderRepository()
	suite.dispatcher = mock.NewMessageDispatcher()

	validator, err := validation.NewValidator()
	if err != nil {
		suite.T().Fatal(err)
	}

	transactionManager := persistence.NilTransactionManager{}
	billing := messaging.NewBillingAdapter(suite.dispatcher)
	tradingAdapter := messaging.NewTradingAdapter(suite.dispatcher)
	dealer := trading.NewDealer(
		suite.items,
		suite.userItems,
		suite.purchaseOrders,
		suite.sellOrders,
		transactionManager,
		billing,
		tradingAdapter,
	)
	service := api.NewTradingApiService(
		suite.purchaseOrders,
		suite.sellOrders,
		suite.items,
		suite.userItems,
		suite.userItems,
		transactionManager,
		dealer,
		validator,
		mock.Locker{},
		0,
	)
	controller := api.NewTradingApiController(service)
	suite.api = api.NewRouter(controller)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func (suite *APISuite) GivenInitialOrder(item *trading.Item, userItem *trading.UserItem) *trading.SellOrder {
	return trading.NewInitialOrder(brokerID, item, userItem)
}

func getPurchaseIdempotenceKey(userID uuid.UUID) string {
	hash := sha256.Sum256([]byte(userID.String() + ":purchase:test"))

	return hex.EncodeToString(hash[:])
}

func getSellIdempotenceKey(userID uuid.UUID) string {
	hash := sha256.Sum256([]byte(userID.String() + ":sell:test"))

	return hex.EncodeToString(hash[:])
}
