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

type UserItemRepository struct {
	conn postgres.Connection
}

func NewUserItemRepository(conn postgres.Connection) *UserItemRepository {
	return &UserItemRepository{conn: conn}
}

func (repository *UserItemRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*trading.AggregatedUserItem, error) {
	rows, err := queries(ctx, repository.conn).FindAggregatedUserItems(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find aggregated user items")
	}

	items := make([]*trading.AggregatedUserItem, len(rows))
	for i, row := range rows {
		items[i] = &trading.AggregatedUserItem{
			ID:          row.ID,
			Name:        row.Name,
			Count:       int(row.Count),
			OnSaleCount: int(row.OnSaleCount),
		}
	}

	return items, nil
}

func (repository *UserItemRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*trading.UserItem, error) {
	row, err := queries(ctx, repository.conn).FindUserItemByIDForUpdate(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrItemNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user item")
	}

	item := userItemFromRow(row)

	return item, nil
}

func (repository *UserItemRepository) FindForSale(ctx context.Context, userID uuid.UUID) (*trading.UserItem, error) {
	row, err := queries(ctx, repository.conn).FindUserItemForSale(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.WithStack(trading.ErrItemNotFound)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user item")
	}

	item := userItemFromRow(row)

	return item, nil
}

func (repository *UserItemRepository) Save(ctx context.Context, item *trading.UserItem) error {
	if item.IsNew() {
		err := queries(ctx, repository.conn).CreateUserItem(ctx, database.CreateUserItemParams{
			ID:       item.ID,
			ItemID:   item.ItemID,
			UserID:   item.UserID,
			IsOnSale: item.IsOnSale,
		})
		if err != nil {
			return errors.Wrap(err, "failed to insert user item")
		}
	} else {
		err := queries(ctx, repository.conn).UpdateUserItem(ctx, database.UpdateUserItemParams{
			ID:       item.ID,
			IsOnSale: item.IsOnSale,
		})
		if err != nil {
			return errors.Wrap(err, "failed to insert user item")
		}
	}

	return nil
}

func (repository *UserItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := queries(ctx, repository.conn).DeleteUserItem(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user item")
	}

	return nil
}

func userItemFromRow(row database.UserItem) *trading.UserItem {
	return &trading.UserItem{
		ID:        row.ID,
		ItemID:    row.ItemID,
		UserID:    row.UserID,
		IsOnSale:  row.IsOnSale,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
