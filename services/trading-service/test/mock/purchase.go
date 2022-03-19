package mock

import (
	"context"
	"testing"

	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
)

type PurchaseOrderRepository struct {
	orders map[uuid.UUID]*trading.PurchaseOrder
}

func NewPurchaseOrderRepository() *PurchaseOrderRepository {
	return &PurchaseOrderRepository{orders: map[uuid.UUID]*trading.PurchaseOrder{}}
}

func (repository *PurchaseOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*trading.PurchaseOrder, error) {
	order := repository.orders[id]
	if order == nil {
		return nil, trading.ErrOrderNotFound
	}

	return order, nil
}

func (repository *PurchaseOrderRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*trading.PurchaseOrder, error) {
	return repository.FindByID(ctx, id)
}

func (repository *PurchaseOrderRepository) FindByPaymentForUpdate(ctx context.Context, paymentID uuid.UUID) (*trading.PurchaseOrder, error) {
	for _, order := range repository.orders {
		if order.PaymentID.UUID == paymentID {
			return order, nil
		}
	}

	return nil, trading.ErrOrderNotFound
}

func (repository *PurchaseOrderRepository) FindByDealForUpdate(ctx context.Context, dealID uuid.UUID) (*trading.PurchaseOrder, error) {
	for _, order := range repository.orders {
		if order.DealID.UUID == dealID {
			return order, nil
		}
	}

	return nil, trading.ErrOrderNotFound
}

func (repository *PurchaseOrderRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*trading.PurchaseOrder, error) {
	orders := make([]*trading.PurchaseOrder, 0)

	for _, order := range repository.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (repository *PurchaseOrderRepository) FindForDeal(ctx context.Context, userID, itemID uuid.UUID, minPrice float64) (*trading.PurchaseOrder, error) {
	var bestOrder *trading.PurchaseOrder

	for _, order := range repository.orders {
		if order.ItemID == itemID &&
			order.UserID != userID &&
			(bestOrder == nil || order.TotalPrice() < bestOrder.TotalPrice()) &&
			order.TotalPrice() >= minPrice &&
			order.Status == trading.PurchasePending {
			bestOrder = order
		}
	}
	if bestOrder == nil {
		return nil, trading.ErrOrderNotFound
	}

	return bestOrder, nil
}

func (repository *PurchaseOrderRepository) Save(ctx context.Context, order *trading.PurchaseOrder) error {
	repository.orders[order.ID] = order

	return nil
}

func (repository *PurchaseOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(repository.orders, id)

	return nil
}

func (repository *PurchaseOrderRepository) Set(orders ...*trading.PurchaseOrder) {
	for _, order := range orders {
		repository.orders[order.ID] = order
	}
}

func (repository *PurchaseOrderRepository) Assert(
	tb testing.TB,
	id uuid.UUID,
	assert func(order *trading.PurchaseOrder),
) {
	order := repository.orders[id]
	if order == nil {
		tb.Errorf("purchase order %s does not exist", id)
	} else {
		assert(order)
	}
}
