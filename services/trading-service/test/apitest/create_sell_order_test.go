package apitest

import (
	"strings"
	"testing"

	"trading-service/internal/messaging"
	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
	"github.com/muonsoft/api-testing/apitest"
	"github.com/muonsoft/api-testing/assertjson"
	"github.com/stretchr/testify/assert"
)

func (suite *APISuite) TestCreateSellOrder_WhenUnauthorized_ExpectUnauthorizedError() {
	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/sell-orders",
		strings.NewReader(`{}`),
	)

	response.IsUnauthorized()
}

func (suite *APISuite) TestCreateSellOrder_WhenInvalidForm_ExpectValidationError() {
	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/sell-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 0
		}`),
		apitest.WithHeader("X-User-Id", userID.String()),
		apitest.WithHeader("If-Match", getPurchaseIdempotenceKey(userID)),
	)

	response.IsUnprocessableEntity()
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/message").EqualToTheString("violation at 'price': This value should be greater than 1.")
	})
}

func (suite *APISuite) TestCreateSellOrder_WhenItemDoesNotExist_ExpectValidationError() {
	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/sell-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 100
		}`),
		apitest.WithHeader("X-User-Id", userID.String()),
		apitest.WithHeader("If-Match", getSellIdempotenceKey(userID)),
	)

	response.IsUnprocessableEntity()
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/message").EqualToTheString("trading item does not exist")
	})
}

func (suite *APISuite) TestCreateSellOrder_WhenNoPurchaseOrderFound_ExpectPendingSellOrderCreated() {
	item := &trading.Item{
		ID:                itemID,
		Name:              "testItem",
		InitialCount:      10,
		InitialPrice:      100,
		CommissionPercent: 10,
	}
	userItem := &trading.UserItem{
		ID:     userItemID,
		ItemID: itemID,
		UserID: sellerID,
	}
	suite.items.Set(item)
	suite.userItems.Set(userItem)

	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/sell-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 100
		}`),
		apitest.WithHeader("X-User-Id", sellerID.String()),
		apitest.WithHeader("If-Match", getSellIdempotenceKey(sellerID)),
	)

	response.IsAccepted()
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/id").AssertString(func(t testing.TB, value string) {
			suite.sellOrders.Assert(t, uuid.FromStringOrNil(value), func(order *trading.SellOrder) {
				assert.Equal(t, trading.SellPending, order.Status)
			})
		})
	})
}

func (suite *APISuite) TestCreateSellOrder_WhenInitialSellOrderFound_ExpectDealInitiation() {
	item := &trading.Item{
		ID:                itemID,
		Name:              "testItem",
		InitialCount:      10,
		InitialPrice:      100,
		CommissionPercent: 10,
	}
	userItem := &trading.UserItem{
		ID:     userItemID,
		ItemID: itemID,
		UserID: sellerID,
	}
	suite.items.Set(item)
	suite.userItems.Set(userItem)
	purchaseOrder := trading.NewPurchaseOrder(purchaserID, item, 111)
	suite.purchaseOrders.Set(
		purchaseOrder,
		trading.NewPurchaseOrder(purchaserID, item, 120),
		trading.NewPurchaseOrder(purchaserID, item, 100),
	)

	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/sell-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 120
		}`),
		apitest.WithHeader("X-User-Id", sellerID.String()),
		apitest.WithHeader("If-Match", getSellIdempotenceKey(sellerID)),
	)

	response.IsAccepted()
	paymentID := uuid.Nil
	dealID := uuid.Nil
	suite.purchaseOrders.Assert(suite.T(), purchaseOrder.ID, func(order *trading.PurchaseOrder) {
		suite.Equal(trading.PurchasePaymentPending, order.Status)
		suite.True(order.PaymentID.Valid)
		suite.True(order.DealID.Valid)
		paymentID = order.PaymentID.UUID
		dealID = order.DealID.UUID
	})
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/id").AssertString(func(t testing.TB, value string) {
			suite.sellOrders.Assert(t, uuid.FromStringOrNil(value), func(order *trading.SellOrder) {
				assert.Equal(t, trading.SellDealPending, order.Status)
				assert.Equal(t, sellerID.String(), order.UserID.String())
				assert.True(t, order.DealID.Valid)
				assert.Equal(t, dealID, order.DealID.UUID)
			})
		})
	})
	suite.dispatcher.AssertMessage(suite.T(), 0, func(t *testing.T, message messaging.Message) {
		payment, ok := message.(messaging.CreatePayment)
		if !ok {
			t.Fatal("expecting messaging.CreatePayment")
		}
		assert.Equal(t, paymentID.String(), payment.ID.String())
		assert.Equal(t, purchaserID.String(), payment.UserID.String())
		assert.Equal(t, 111.0, payment.Amount)
		assert.Equal(t, 11.1, payment.Commission)
		assert.Equal(t, `buying the item "testItem" on the marketplace (with commission 11)`, payment.Description)
	})
}
