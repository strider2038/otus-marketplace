package redis

import (
	"context"
	"time"

	"trading-service/internal/api"

	"github.com/bsm/redislock"
)

type LockerAdapter struct {
	locker *redislock.Client
}

func NewLockerAdapter(locker *redislock.Client) *LockerAdapter {
	return &LockerAdapter{locker: locker}
}

func (locker *LockerAdapter) Obtain(ctx context.Context, key string, ttl time.Duration) (api.Lock, error) {
	return locker.locker.Obtain(ctx, key, ttl, nil)
}
