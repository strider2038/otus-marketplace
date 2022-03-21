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

func (suite *APISuite) TestCreatePurchaseOrder_WhenUnauthorized_ExpectUnauthorizedError() {
	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/purchase-orders",
		strings.NewReader(`{}`),
	)

	response.IsUnauthorized()
}

func (suite *APISuite) TestCreatePurchaseOrder_WhenInvalidForm_ExpectValidationError() {
	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/purchase-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 0
		}`),
		apitest.WithHeader("X-User-Id", userID.String()),
	)

	response.IsUnprocessableEntity()
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/message").EqualToTheString("violation at 'price': This value should be greater than 1.")
	})
}

func (suite *APISuite) TestCreatePurchaseOrder_WhenItemDoesNotExist_ExpectValidationError() {
	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/purchase-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 100
		}`),
		apitest.WithHeader("X-User-Id", userID.String()),
	)

	response.IsUnprocessableEntity()
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/message").EqualToTheString("trading item does not exist")
	})
}

func (suite *APISuite) TestCreatePurchaseOrder_WhenItemNoSellOrderFound_ExpectPendingPurchaseOrderCreated() {
	suite.items.Set(
		&trading.Item{
			ID:                itemID,
			Name:              "testItem",
			InitialCount:      10,
			InitialPrice:      100,
			CommissionPercent: 10,
		},
	)

	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/purchase-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 100
		}`),
		apitest.WithHeader("X-User-Id", userID.String()),
	)

	response.IsAccepted()
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/id").AssertString(func(t testing.TB, value string) {
			suite.purchaseOrders.Assert(t, uuid.FromStringOrNil(value), func(order *trading.PurchaseOrder) {
				assert.Equal(t, trading.PurchasePending, order.Status)
			})
		})
	})
}

func (suite *APISuite) TestCreatePurchaseOrder_WhenInitialSellOrderFound_ExpectDealInitiation() {
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
		UserID: brokerID,
	}
	suite.items.Set(item)
	suite.userItems.Set(userItem)
	sellOrder := suite.GivenInitialOrder(item, userItem)
	suite.sellOrders.Set(sellOrder)

	response := apitest.HandlePOST(
		suite.T(),
		suite.api,
		"/api/v1/purchase-orders",
		strings.NewReader(`{
			"itemId": "31e2951a-08da-4ae6-b37b-928cf2507448",
			"price": 120
		}`),
		apitest.WithHeader("X-User-Id", purchaserID.String()),
	)

	response.IsAccepted()
	paymentID := uuid.Nil
	dealID := uuid.Nil
	response.HasJSON(func(json *assertjson.AssertJSON) {
		json.Node("/id").AssertString(func(t testing.TB, value string) {
			suite.purchaseOrders.Assert(t, uuid.FromStringOrNil(value), func(order *trading.PurchaseOrder) {
				assert.Equal(t, trading.PurchasePaymentPending, order.Status)
				assert.Equal(t, purchaserID.String(), order.UserID.String())
				assert.True(t, order.PaymentID.Valid)
				assert.True(t, order.DealID.Valid)
				paymentID = order.PaymentID.UUID
				dealID = order.DealID.UUID
			})
		})
	})
	suite.sellOrders.Assert(suite.T(), sellOrder.ID, func(order *trading.SellOrder) {
		suite.Equal(trading.SellDealPending, order.Status)
		suite.True(order.DealID.Valid)
		suite.Equal(dealID, order.DealID.UUID)
	})
	suite.dispatcher.AssertMessage(suite.T(), 0, func(t *testing.T, message messaging.Message) {
		payment, ok := message.(messaging.CreatePayment)
		if !ok {
			t.Fatal("expecting messaging.CreatePayment")
		}
		assert.Equal(t, paymentID.String(), payment.ID.String())
		assert.Equal(t, purchaserID.String(), payment.UserID.String())
		assert.Equal(t, 100.0, payment.Amount)
		assert.Equal(t, 10.0, payment.Commission)
		assert.Equal(t, `buying the item "testItem" on the marketplace (with commission 10)`, payment.Description)
	})
}
