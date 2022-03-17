// Code generated by sqlc. DO NOT EDIT.

package database

import (
	"time"

	"github.com/gofrs/uuid"
)

type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Message   string
	CreatedAt time.Time
}

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
