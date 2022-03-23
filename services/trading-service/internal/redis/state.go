package redis

import (
	"context"
	"time"

	"trading-service/internal/trading"

	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type StateRepository struct {
	namespace  string
	redis      *redis.Client
	expiration time.Duration
}

func NewStateRepository(namespace string, redis *redis.Client, expiration time.Duration) *StateRepository {
	return &StateRepository{namespace: namespace, redis: redis, expiration: expiration}
}

func (s *StateRepository) Get(ctx context.Context, userID uuid.UUID) (string, error) {
	key := s.getKey(userID)
	value, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return uuid.Nil.String(), nil
	}
	if err != nil {
		return "", errors.Wrapf(err, `failed to get state by key "%s"`, key)
	}

	return value, nil
}

func (s *StateRepository) Verify(ctx context.Context, userID uuid.UUID, value string) error {
	expected, err := s.Get(ctx, userID)
	if err != nil {
		return err
	}
	if expected != value {
		return errors.WithStack(trading.ErrOutdated)
	}

	return nil
}

func (s *StateRepository) Refresh(ctx context.Context, userID uuid.UUID) error {
	key := s.getKey(userID)
	value := uuid.Must(uuid.NewV4()).String()
	err := s.redis.Set(ctx, key, value, s.expiration).Err()
	if err != nil {
		return errors.Wrapf(err, `failed to save new state into redis by key "%s"`, key)
	}

	return nil
}

func (s *StateRepository) getKey(userID uuid.UUID) string {
	return s.namespace + ":" + userID.String()
}
