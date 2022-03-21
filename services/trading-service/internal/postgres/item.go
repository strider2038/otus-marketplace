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

type ItemRepository struct {
	conn postgres.Connection
}

func NewItemRepository(conn postgres.Connection) *ItemRepository {
	return &ItemRepository{conn: conn}
}

func (repository *ItemRepository) FindAll(ctx context.Context) ([]*trading.Item, error) {
	rows, err := queries(ctx, repository.conn).FindAllItems(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find all items")
	}

	items := make([]*trading.Item, len(rows))
	for i, row := range rows {
		items[i] = itemFromRow(row)
	}

	return items, nil
}

func (repository *ItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*trading.Item, error) {
	row, err := queries(ctx, repository.conn).FindItem(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrItemNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find item")
	}

	return itemFromRow(row), nil
}

func (repository *ItemRepository) Add(ctx context.Context, item *trading.Item) error {
	err := queries(ctx, repository.conn).AddItem(ctx, database.AddItemParams{
		ID:                item.ID,
		Name:              item.Name,
		InitialCount:      item.InitialCount,
		InitialPrice:      item.InitialPrice,
		CommissionPercent: item.CommissionPercent,
		CreatedAt:         item.CreatedAt,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add item")
	}

	return nil
}

func itemFromRow(row database.Item) *trading.Item {
	return &trading.Item{
		ID:                row.ID,
		Name:              row.Name,
		InitialCount:      row.InitialCount,
		InitialPrice:      row.InitialPrice,
		CommissionPercent: row.CommissionPercent,
		CreatedAt:         row.CreatedAt,
	}
}
