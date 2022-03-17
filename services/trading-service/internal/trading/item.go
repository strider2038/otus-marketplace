package trading

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Item struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	InitialCount      int32     `json:"initialCount"`
	InitialPrice      float64   `json:"initialPrice"`
	CommissionPercent float64   `json:"commissionPercent"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (item Item) CalculateCommission(price float64) float64 {
	commission := price * item.CommissionPercent / 100

	return commission
}

func NewItem(name string, initialCount int32, initialPrice, commissionPercent float64) *Item {
	return &Item{
		ID:                uuid.Must(uuid.NewV4()),
		Name:              name,
		InitialCount:      initialCount,
		InitialPrice:      initialPrice,
		CommissionPercent: commissionPercent,
		CreatedAt:         time.Now(),
	}
}

type ItemRepository interface {
	FindAll(ctx context.Context) ([]*Item, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Item, error)
	Add(ctx context.Context, item *Item) error
}
