package statistics

import (
	"context"

	"github.com/gofrs/uuid"
)

type DailyDeals struct {
	Date     string    `json:"date"`
	ItemID   uuid.UUID `json:"itemId"`
	ItemName string    `json:"itemName"`
	Count    int64     `json:"count"`
	Amount   float64   `json:"amount"`
}

type DailyDealsRepository interface {
	FindForLastWeek(ctx context.Context) ([]*DailyDeals, error)
	Add(ctx context.Context, deals *DailyDeals) error
}
