package statistics

import "context"

type TotalDailyDeals struct {
	Date   string  `json:"date"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

type TotalDailyDealsRepository interface {
	FindForLastWeek(ctx context.Context) ([]*TotalDailyDeals, error)
	Add(ctx context.Context, deals *TotalDailyDeals) error
}
