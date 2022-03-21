package postgres

import (
	"context"

	"history-service/internal/history"
	"history-service/internal/postgres/database"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type DealRepository struct {
	db *database.Queries
}

func NewDealRepository(db *database.Queries) *DealRepository {
	return &DealRepository{db: db}
}

func (repository *DealRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*history.Deal, error) {
	rows, err := repository.db.FindDealsByUser(ctx, userID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to find deals by user")
	}

	deals := make([]*history.Deal, len(rows))
	for i, row := range rows {
		deals[i] = &history.Deal{
			ID:          row.ID,
			UserID:      row.UserID,
			ItemID:      row.ItemID,
			ItemName:    row.ItemName,
			Type:        history.DealType(row.Type),
			Amount:      row.Amount,
			Commission:  row.Commission,
			CompletedAt: row.CompletedAt,
		}
	}

	return deals, nil
}

func (repository *DealRepository) Add(ctx context.Context, deal *history.Deal) error {
	err := repository.db.AddDeal(ctx, database.AddDealParams{
		ID:          deal.ID,
		UserID:      deal.UserID,
		ItemID:      deal.ItemID,
		ItemName:    deal.ItemName,
		Type:        database.DealType(deal.Type),
		Amount:      deal.Amount,
		Commission:  deal.Commission,
		CompletedAt: deal.CompletedAt,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to add deal")
	}

	return nil
}
