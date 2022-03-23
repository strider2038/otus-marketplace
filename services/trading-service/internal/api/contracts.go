package api

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Locker interface {
	Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error)
}

type Lock interface {
	Release(ctx context.Context) error
}

type StateRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (string, error)
	Verify(ctx context.Context, userID uuid.UUID, value string) error
	Refresh(ctx context.Context, userID uuid.UUID) error
}
