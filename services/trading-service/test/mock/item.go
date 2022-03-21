package mock

import (
	"context"

	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
)

type ItemRepository struct {
	items map[uuid.UUID]*trading.Item
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{items: map[uuid.UUID]*trading.Item{}}
}

func (repository *ItemRepository) FindAll(ctx context.Context) ([]*trading.Item, error) {
	items := make([]*trading.Item, 0, len(repository.items))
	for _, item := range repository.items {
		items = append(items, item)
	}

	return items, nil
}

func (repository *ItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*trading.Item, error) {
	item := repository.items[id]
	if item == nil {
		return nil, trading.ErrItemNotFound
	}

	return item, nil
}

func (repository *ItemRepository) Add(ctx context.Context, item *trading.Item) error {
	repository.items[item.ID] = item

	return nil
}

func (repository *ItemRepository) Set(items ...*trading.Item) {
	for _, item := range items {
		repository.items[item.ID] = item
	}
}
