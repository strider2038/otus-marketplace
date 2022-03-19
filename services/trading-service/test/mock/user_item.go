package mock

import (
	"context"
	"fmt"

	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
)

type UserItemRepository struct {
	items map[uuid.UUID]*trading.UserItem
}

func NewUserItemRepository() *UserItemRepository {
	return &UserItemRepository{items: map[uuid.UUID]*trading.UserItem{}}
}

func (repository *UserItemRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*trading.AggregatedUserItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (repository *UserItemRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*trading.UserItem, error) {
	item := repository.items[id]
	if item == nil {
		return nil, trading.ErrItemNotFound
	}

	return item, nil
}

func (repository *UserItemRepository) FindForSale(ctx context.Context, userID uuid.UUID) (*trading.UserItem, error) {
	for _, item := range repository.items {
		if item.UserID == userID && !item.IsOnSale {
			return item, nil
		}
	}

	return nil, trading.ErrItemNotFound
}

func (repository *UserItemRepository) Save(ctx context.Context, item *trading.UserItem) error {
	repository.items[item.ID] = item

	return nil
}

func (repository *UserItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(repository.items, id)

	return nil
}

func (repository *UserItemRepository) Set(items ...*trading.UserItem) {
	for _, item := range items {
		repository.items[item.ID] = item
	}
}
