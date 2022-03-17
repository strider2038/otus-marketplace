package api

import (
	"context"

	"identity-service/internal/events"
	"identity-service/internal/users"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type TokenIssuer interface {
	Issue(user *users.User) (string, error)
}

type EventDispatcher interface {
	Dispatch(ctx context.Context, event events.Event) error
}
