package trading

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type UserItem struct {
	isNew bool

	ID        uuid.UUID
	ItemID    uuid.UUID
	UserID    uuid.UUID
	IsOnSale  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u UserItem) IsNew() bool { return u.isNew }

func NewUserItem(itemID uuid.UUID, userID uuid.UUID) *UserItem {
	now := time.Now()

	return &UserItem{
		isNew:     true,
		ID:        uuid.Must(uuid.NewV4()),
		ItemID:    itemID,
		UserID:    userID,
		IsOnSale:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type AggregatedUserItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Count       int       `json:"count"`
	OnSaleCount int       `json:"onSaleCount"`
}

type UserItemRepository interface {
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*UserItem, error)
	FindForSale(ctx context.Context, userID uuid.UUID) (*UserItem, error)
	Save(ctx context.Context, item *UserItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AggregatedUserItemRepository interface {
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*AggregatedUserItem, error)
}
