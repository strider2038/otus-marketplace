package mock

import (
	"context"

	"github.com/gofrs/uuid"
)

type StateRepository struct{}

func (s StateRepository) Get(ctx context.Context, userID uuid.UUID) (string, error) {
	return "test-state", nil
}

func (s StateRepository) Verify(ctx context.Context, userID uuid.UUID, value string) error {
	return nil
}

func (s StateRepository) Refresh(ctx context.Context, userID uuid.UUID) error {
	return nil
}
