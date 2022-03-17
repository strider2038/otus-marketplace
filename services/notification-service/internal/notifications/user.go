package notifications

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	Save(ctx context.Context, user *User) error
}
