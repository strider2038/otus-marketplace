package postgres

import (
	"context"
	"errors"
	"fmt"

	"identity-service/internal/postgres/database"
	"identity-service/internal/users"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type UserRepository struct {
	db *database.Queries
}

func NewUserRepository(db *database.Queries) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*users.User, error) {
	row, err := repository.db.FindUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, users.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return userFromRow(row), nil
}

func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	row, err := repository.db.FindUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, users.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return userFromRow(row), nil
}

func (repository *UserRepository) CountByEmail(ctx context.Context, email string) (int64, error) {
	count, err := repository.db.CountUsersByEmail(ctx, email)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by email: %w", err)
	}
	return count, nil
}

func (repository *UserRepository) Save(ctx context.Context, user *users.User) error {
	var u database.User
	var err error

	if user.ID.IsNil() {
		u, err = repository.db.CreateUser(ctx, database.CreateUserParams{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     user.Email,
			Password:  user.Password,
			Role:      database.UserRole(user.Role),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		user.ID = u.ID
	} else {
		u, err = repository.db.UpdateUser(ctx, database.UpdateUserParams{
			ID:        user.ID,
			Email:     user.Email,
			Role:      database.UserRole(user.Role),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		})
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	user.Email = u.Email
	user.Password = u.Password
	user.Role = users.Role(u.Role)
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.Phone = u.Phone
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt

	return nil
}

func (repository *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := repository.db.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func userFromRow(row database.User) *users.User {
	return &users.User{
		ID:        row.ID,
		Email:     row.Email,
		Role:      users.Role(row.Role),
		Password:  row.Password,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Phone:     row.Phone,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
