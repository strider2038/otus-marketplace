package postgres

import (
	"context"

	notificationspkg "notification-service/internal/notifications"
	"notification-service/internal/postgres/database"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *database.Queries
}

func NewUserRepository(db *database.Queries) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*notificationspkg.User, error) {
	user, err := repository.db.FindUser(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}

	return userFromDB(user), nil
}

func (repository *UserRepository) Save(ctx context.Context, user *notificationspkg.User) error {
	_, err := repository.db.CreateUser(ctx, database.CreateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func userFromDB(user database.User) *notificationspkg.User {
	return &notificationspkg.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
