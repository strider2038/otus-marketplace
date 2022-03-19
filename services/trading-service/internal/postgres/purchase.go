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

type PurchaseOrderRepository struct {
	conn postgres.Connection
}

func NewPurchaseOrderRepository(conn postgres.Connection) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{conn: conn}
}

func (repository *PurchaseOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*trading.PurchaseOrder, error) {
	row, err := queries(ctx, repository.conn).FindPurchaseOrder(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find purchase order")
	}

	return purchaseOrderFromRow(row), nil
}

func (repository *PurchaseOrderRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*trading.PurchaseOrder, error) {
	row, err := queries(ctx, repository.conn).FindPurchaseOrderForUpdate(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find purchase order")
	}

	return purchaseOrderFromRow(row), nil
}

func (repository *PurchaseOrderRepository) FindByPaymentForUpdate(ctx context.Context, paymentID uuid.UUID) (*trading.PurchaseOrder, error) {
	row, err := queries(ctx, repository.conn).FindPurchaseOrderByPaymentForUpdate(ctx, uuid.NullUUID{UUID: paymentID, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find purchase order")
	}

	return purchaseOrderFromRow(row), nil
}

func (repository *PurchaseOrderRepository) FindByDealForUpdate(ctx context.Context, dealID uuid.UUID) (*trading.PurchaseOrder, error) {
	row, err := queries(ctx, repository.conn).FindPurchaseOrderByDealForUpdate(ctx, uuid.NullUUID{UUID: dealID, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find purchase order")
	}

	return purchaseOrderFromRow(row), nil
}

func (repository *PurchaseOrderRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*trading.PurchaseOrder, error) {
	rows, err := queries(ctx, repository.conn).FindPurchaseOrdersByUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find purchase orders")
	}

	orders := make([]*trading.PurchaseOrder, len(rows))
	for i, row := range rows {
		orders[i] = purchaseOrderFromRow(row)
	}

	return orders, nil
}

func (repository *PurchaseOrderRepository) FindForDeal(ctx context.Context, userID, itemID uuid.UUID, minPrice float64) (*trading.PurchaseOrder, error) {
	row, err := queries(ctx, repository.conn).FindPurchaseOrderForDeal(ctx, database.FindPurchaseOrderForDealParams{
		UserID: userID,
		ItemID: itemID,
		Price:  minPrice,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrOrderNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find purchase order")
	}

	return purchaseOrderFromRow(row), nil
}

func (repository *PurchaseOrderRepository) Save(ctx context.Context, order *trading.PurchaseOrder) error {
	if order.IsNew() {
		return repository.insert(ctx, order)
	}

	return repository.update(ctx, order)
}

func (repository *PurchaseOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := queries(ctx, repository.conn).DeletePurchaseOrder(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete purchase order")
	}

	return nil
}

func (repository *PurchaseOrderRepository) insert(ctx context.Context, order *trading.PurchaseOrder) error {
	err := queries(ctx, repository.conn).CreatePurchaseOrder(ctx, database.CreatePurchaseOrderParams{
		ID:         order.ID,
		UserID:     order.UserID,
		ItemID:     order.ItemID,
		PaymentID:  order.PaymentID,
		DealID:     order.DealID,
		Price:      order.Price,
		Commission: order.Commission,
		Status:     database.PurchaseStatus(order.Status),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	})
	if err != nil {
		return errors.Wrap(err, "failed to insert purchase order")
	}

	return nil
}

func (repository *PurchaseOrderRepository) update(ctx context.Context, order *trading.PurchaseOrder) error {
	err := queries(ctx, repository.conn).UpdatePurchaseOrder(ctx, database.UpdatePurchaseOrderParams{
		ID:        order.ID,
		PaymentID: order.PaymentID,
		DealID:    order.DealID,
		Status:    database.PurchaseStatus(order.Status),
	})
	if err != nil {
		return errors.Wrap(err, "failed to update purchase order")
	}

	return nil
}

func purchaseOrderFromRow(row database.PurchaseOrder) *trading.PurchaseOrder {
	return &trading.PurchaseOrder{
		ID:         row.ID,
		UserID:     row.UserID,
		ItemID:     row.ItemID,
		PaymentID:  row.PaymentID,
		DealID:     row.DealID,
		Price:      row.Price,
		Commission: row.Commission,
		Status:     trading.PurchaseStatus(row.Status),
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}
