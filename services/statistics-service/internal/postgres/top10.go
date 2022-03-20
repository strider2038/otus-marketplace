package postgres

import (
	"context"

	"statistics-service/internal/postgres/database"
	"statistics-service/internal/statistics"

	"github.com/pkg/errors"
)

type Top10DealsRepository struct {
	db *database.Queries
}

func NewTop10DealsRepository(db *database.Queries) *Top10DealsRepository {
	return &Top10DealsRepository{db: db}
}

func (repository *Top10DealsRepository) FindTop10(ctx context.Context) ([]*statistics.Top10Deals, error) {
	rows, err := repository.db.FindTop10Deals(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find top 10 deals")
	}

	deals := make([]*statistics.Top10Deals, len(rows))
	for i, row := range rows {
		deals[i] = &statistics.Top10Deals{
			ItemID:   row.ItemID,
			ItemName: row.ItemName,
			Count:    row.Count,
			Amount:   row.Amount,
		}
	}

	return deals, nil
}

func (repository *Top10DealsRepository) Add(ctx context.Context, deals *statistics.Top10Deals) error {
	err := repository.db.AddTop10Deals(ctx, database.AddTop10DealsParams{
		ItemID:   deals.ItemID,
		ItemName: deals.ItemName,
		Count:    deals.Count,
		Amount:   deals.Amount,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add top 10 deals")
	}

	return nil
}
