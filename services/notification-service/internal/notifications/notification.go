package notifications

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"-"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewNotification(userID uuid.UUID, message string) *Notification {
	return &Notification{
		ID:      uuid.Must(uuid.NewV4()),
		UserID:  userID,
		Message: message,
	}
}

type Repository interface {
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*Notification, error)
	Add(ctx context.Context, notification *Notification) error
}
