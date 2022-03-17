package messaging

import (
	"context"
	"encoding/json"

	"notification-service/internal/notifications"

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
	users notifications.UserRepository
}

func NewUserCreatedProcessor(users notifications.UserRepository) *UserCreatedProcessor {
	return &UserCreatedProcessor{users: users}
}

func (processor *UserCreatedProcessor) Process(ctx context.Context, message []byte) error {
	var user UserCreated
	err := json.Unmarshal(message, &user)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal UserCreated event")
	}

	err = processor.users.Save(ctx, &notifications.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	})
	if err != nil {
		return errors.WithMessagef(err, `failed to create user "%s"`, user.ID)
	}

	return nil
}
