package postgres

import (
	"context"

	"billing-service/internal/billing"
	"billing-service/internal/postgres/database"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	postgres "github.com/strider2038/pkg/persistence/pgx"
)

type OperationRepository struct {
	conn postgres.Connection
}

func NewOperationRepository(conn postgres.Connection) *OperationRepository {
	return &OperationRepository{conn: conn}
}

func (repository *OperationRepository) FindByAccount(ctx context.Context, accountID uuid.UUID) ([]*billing.Operation, error) {
	rows, err := queries(ctx, repository.conn).FindOperationsByAccount(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find operations")
	}

	operations := make([]*billing.Operation, len(rows))
	for i, row := range rows {
		operations[i] = &billing.Operation{
			ID:          row.ID,
			AccountID:   row.AccountID,
			Type:        billing.OperationType(row.Type),
			Amount:      row.Amount,
			Description: row.Description,
			CreatedAt:   row.CreatedAt,
		}
	}

	return operations, nil
}

func (repository *OperationRepository) Add(ctx context.Context, operation *billing.Operation) error {
	_, err := queries(ctx, repository.conn).CreateOperation(ctx, database.CreateOperationParams{
		ID:          operation.ID,
		AccountID:   operation.AccountID,
		Amount:      operation.Amount,
		Type:        database.OperationType(operation.Type),
		Description: operation.Description,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add operation")
	}

	return nil
}
