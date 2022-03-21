package statistics

import (
	"context"

	"github.com/gofrs/uuid"
)

type Top10Deals struct {
	ItemID   uuid.UUID `json:"itemId"`
	ItemName string    `json:"itemName"`
	Count    int64     `json:"count"`
	Amount   float64   `json:"amount"`
}

type Top10DealsRepository interface {
	FindTop10(ctx context.Context) ([]*Top10Deals, error)
	Add(ctx context.Context, deals *Top10Deals) error
}
