package mock

import (
	"context"
	"testing"

	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
)

type SellOrderRepository struct {
	orders map[uuid.UUID]*trading.SellOrder
}

func NewSellOrderRepository() *SellOrderRepository {
	return &SellOrderRepository{orders: map[uuid.UUID]*trading.SellOrder{}}
}

func (repository *SellOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*trading.SellOrder, error) {
	order := repository.orders[id]
	if order == nil {
		return nil, trading.ErrOrderNotFound
	}

	return order, nil
}

func (repository *SellOrderRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*trading.SellOrder, error) {
	return repository.FindByID(ctx, id)
}

func (repository *SellOrderRepository) FindByAccrualForUpdate(ctx context.Context, accrualID uuid.UUID) (*trading.SellOrder, error) {
	for _, order := range repository.orders {
		if order.AccrualID.UUID == accrualID {
			return order, nil
		}
	}

	return nil, trading.ErrOrderNotFound
}

func (repository *SellOrderRepository) FindByDealForUpdate(ctx context.Context, dealID uuid.UUID) (*trading.SellOrder, error) {
	for _, order := range repository.orders {
		if order.DealID.UUID == dealID {
			return order, nil
		}
	}

	return nil, trading.ErrOrderNotFound
}

func (repository *SellOrderRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*trading.SellOrder, error) {
	orders := make([]*trading.SellOrder, 0)

	for _, order := range repository.orders {
		if order.UserID.UUID == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (repository *SellOrderRepository) FindForDeal(ctx context.Context, itemID uuid.UUID, maxPrice float64) (*trading.SellOrder, error) {
	var bestOrder *trading.SellOrder
	var bestPrice float64

	for _, order := range repository.orders {
		if order.ItemID == itemID &&
			order.TotalPrice() > bestPrice &&
			order.TotalPrice() <= maxPrice &&
			order.Status == trading.SellPending {
			bestOrder = order
			bestPrice = order.TotalPrice()
		}
	}
	if bestOrder == nil {
		return nil, trading.ErrOrderNotFound
	}

	return bestOrder, nil
}

func (repository *SellOrderRepository) Save(ctx context.Context, order *trading.SellOrder) error {
	repository.orders[order.ID] = order

	return nil
}

func (repository *SellOrderRepository) Set(orders ...*trading.SellOrder) {
	for _, order := range orders {
		repository.orders[order.ID] = order
	}
}

func (repository *SellOrderRepository) Assert(
	tb testing.TB,
	id uuid.UUID,
	assert func(order *trading.SellOrder),
) {
	order := repository.orders[id]
	if order == nil {
		tb.Errorf("sell order %s does not exist", id)
	} else {
		assert(order)
	}
}
