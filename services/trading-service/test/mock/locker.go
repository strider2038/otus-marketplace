package mock

import (
	"context"
	"time"

	"trading-service/internal/api"
)

type Locker struct{}

func (mock Locker) Obtain(ctx context.Context, key string, ttl time.Duration) (api.Lock, error) {
	return Lock{}, nil
}

type Lock struct{}

func (mock Lock) Release(ctx context.Context) error {
	return nil
}
