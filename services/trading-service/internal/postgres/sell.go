package postgres

import (
	"context"

	"trading-service/internal/postgres/database"
	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	postgres "github.com/strider2038/pkg/persistence/pgx"
)

type SellOrderRepository struct {
	conn postgres.Connection
}

func NewSellOrderRepository(conn postgres.Connection) *SellOrderRepository {
	return &SellOrderRepository{conn: conn}
}

func (repository *SellOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*trading.SellOrder, error) {
	row, err := queries(ctx, repository.conn).FindSellOrder(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sell order")
	}

	return sellOrderFromRow(row), nil
}

func (repository *SellOrderRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*trading.SellOrder, error) {
	row, err := queries(ctx, repository.conn).FindSellOrderForUpdate(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sell order")
	}

	return sellOrderFromRow(row), nil
}

func (repository *SellOrderRepository) FindByAccrualForUpdate(ctx context.Context, accrualID uuid.UUID) (*trading.SellOrder, error) {
	row, err := queries(ctx, repository.conn).FindSellOrderByAccrualForUpdate(ctx, uuid.NullUUID{UUID: accrualID, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sell order")
	}

	return sellOrderFromRow(row), nil
}

func (repository *SellOrderRepository) FindByDealForUpdate(ctx context.Context, dealID uuid.UUID) (*trading.SellOrder, error) {
	row, err := queries(ctx, repository.conn).FindSellOrderByDealForUpdate(ctx, uuid.NullUUID{UUID: dealID, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sell order")
	}

	return sellOrderFromRow(row), nil
}

func (repository *SellOrderRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*trading.SellOrder, error) {
	rows, err := queries(ctx, repository.conn).FindSellOrdersByUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sell orders")
	}

	orders := make([]*trading.SellOrder, len(rows))
	for i, row := range rows {
		orders[i] = sellOrderFromRow(row)
	}

	return orders, nil
}

func (repository *SellOrderRepository) FindForDeal(ctx context.Context, userID, itemID uuid.UUID, maxPrice float64) (*trading.SellOrder, error) {
	row, err := queries(ctx, repository.conn).FindSellOrderForDeal(ctx, database.FindSellOrderForDealParams{
		ItemID: itemID,
		UserID: userID,
		Price:  maxPrice,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find sell order")
	}

	return sellOrderFromRow(row), nil
}

func (repository *SellOrderRepository) Save(ctx context.Context, order *trading.SellOrder) error {
	if order.IsNew() {
		return repository.insert(ctx, order)
	}

	return repository.update(ctx, order)
}

func (repository *SellOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := queries(ctx, repository.conn).DeleteSellOrder(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete sell order")
	}

	return nil
}

func (repository *SellOrderRepository) insert(ctx context.Context, order *trading.SellOrder) error {
	err := queries(ctx, repository.conn).CreateSellOrder(ctx, database.CreateSellOrderParams{
		ID:         order.ID,
		UserID:     order.UserID,
		ItemID:     order.ItemID,
		UserItemID: order.UserItemID,
		AccrualID:  order.AccrualID,
		DealID:     order.DealID,
		Price:      order.Price,
		Commission: order.Commission,
		Status:     database.SellStatus(order.Status),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	})
	if err != nil {
		return errors.Wrap(err, "failed to insert sell order")
	}

	return nil
}

func (repository *SellOrderRepository) update(ctx context.Context, order *trading.SellOrder) error {
	err := queries(ctx, repository.conn).UpdateSellOrder(ctx, database.UpdateSellOrderParams{
		ID:        order.ID,
		AccrualID: order.AccrualID,
		DealID:    order.DealID,
		Status:    database.SellStatus(order.Status),
	})
	if err != nil {
		return errors.Wrap(err, "failed to update sell order")
	}

	return nil
}

func (repository *SellOrderRepository) GetStateByUser(ctx context.Context, userID uuid.UUID) (string, error) {
	state, err := queries(ctx, repository.conn).GetSellOrdersStateOfUser(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get sell orders state")
	}

	return state, nil
}

func sellOrderFromRow(row database.SellOrder) *trading.SellOrder {
	return &trading.SellOrder{
		ID:         row.ID,
		ItemID:     row.ItemID,
		UserItemID: row.UserItemID,
		UserID:     row.UserID,
		AccrualID:  row.AccrualID,
		DealID:     row.DealID,
		Price:      row.Price,
		Commission: row.Commission,
		Status:     trading.SellStatus(row.Status),
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
