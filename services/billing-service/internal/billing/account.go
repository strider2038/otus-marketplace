package billing

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var ErrAccountNotFound = errors.New("account not found")

type Account struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type AccountRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Account, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	Create(ctx context.Context, id uuid.UUID) (*Account, error)
	Save(ctx context.Context, account *Account) error
}
