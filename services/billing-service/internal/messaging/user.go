package messaging

import (
	"context"
	"encoding/json"

	"billing-service/internal/billing"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type UserCreated struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Phone     string    `json:"phone,omitempty"`
}

type UserCreatedProcessor struct {
	accounts billing.AccountRepository
}

func NewUserCreatedProcessor(accounts billing.AccountRepository) *UserCreatedProcessor {
	return &UserCreatedProcessor{accounts: accounts}
}

func (processor *UserCreatedProcessor) Process(ctx context.Context, message []byte) error {
	var user UserCreated
	err := json.Unmarshal(message, &user)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal UserCreated event")
	}

	_, err = processor.accounts.Create(ctx, user.ID)
	if err != nil {
		return errors.WithMessagef(err, `failed to create billing account for user "%s"`, user.ID)
	}

	return nil
}
