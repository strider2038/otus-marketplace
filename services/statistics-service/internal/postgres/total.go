package postgres

import (
	"context"
	"time"

	"statistics-service/internal/postgres/database"
	"statistics-service/internal/statistics"

	"github.com/pkg/errors"
)

type TotalDailyDealsRepository struct {
	db *database.Queries
}

func NewTotalDailyDealsRepository(db *database.Queries) *TotalDailyDealsRepository {
	return &TotalDailyDealsRepository{db: db}
}

func (repository *TotalDailyDealsRepository) FindForLastWeek(ctx context.Context) ([]*statistics.TotalDailyDeals, error) {
	rows, err := repository.db.FindTotalDailyDeals(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find daily deals")
	}

	deals := make([]*statistics.TotalDailyDeals, len(rows))
	for i, row := range rows {
		deals[i] = &statistics.TotalDailyDeals{
			Date:   row.Date.Format(statistics.DateLayout),
			Count:  row.Count,
			Amount: row.Amount,
		}
	}

	return deals, nil
}

func (repository *TotalDailyDealsRepository) Add(ctx context.Context, deals *statistics.TotalDailyDeals) error {
	date, err := time.Parse(statistics.DateLayout, deals.Date)
	if err != nil {
		return errors.Wrapf(err, `failed to parse date "%s"`, deals.Date)
	}

	err = repository.db.AddTotalDailyDeals(ctx, database.AddTotalDailyDealsParams{
		Date:   date,
		Count:  deals.Count,
		Amount: deals.Amount,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add total daily deals")
	}

	return nil
}
