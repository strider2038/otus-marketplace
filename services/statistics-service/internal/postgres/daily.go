package postgres

import (
	"context"
	"time"

	"statistics-service/internal/postgres/database"
	"statistics-service/internal/statistics"

	"github.com/pkg/errors"
)

type DailyDealsRepository struct {
	db *database.Queries
}

func NewDailyDealsRepository(db *database.Queries) *DailyDealsRepository {
	return &DailyDealsRepository{db: db}
}

func (repository *DailyDealsRepository) FindForLastWeek(ctx context.Context) ([]*statistics.DailyDeals, error) {
	rows, err := repository.db.FindDailyDeals(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find daily deals")
	}

	deals := make([]*statistics.DailyDeals, len(rows))
	for i, row := range rows {
		deals[i] = &statistics.DailyDeals{
			Date:     row.Date.Format(statistics.DateLayout),
			ItemID:   row.ItemID,
			ItemName: row.ItemName,
			Count:    row.Count,
			Amount:   row.Amount,
		}
	}

	return deals, nil
}

func (repository *DailyDealsRepository) Add(ctx context.Context, deals *statistics.DailyDeals) error {
	date, err := time.Parse(statistics.DateLayout, deals.Date)
	if err != nil {
		return errors.Wrapf(err, `failed to parse date "%s"`, deals.Date)
	}

	err = repository.db.AddDailyDeals(ctx, database.AddDailyDealsParams{
		Date:     date,
		ItemID:   deals.ItemID,
		ItemName: deals.ItemName,
		Count:    deals.Count,
		Amount:   deals.Amount,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add daily deals")
	}

	return nil
}
