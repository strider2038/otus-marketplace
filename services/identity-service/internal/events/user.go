package events

import "github.com/gofrs/uuid"

type UserCreated struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Phone     string    `json:"phone,omitempty"`
}

func (u UserCreated) Name() string {
	return "Identity/UserCreated"
}

type UserUpdated struct {
	ID        uuid.UUID `json:"id,omitempty"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Phone     string    `json:"phone,omitempty"`
}

func (u UserUpdated) Name() string {
	return "Identity/UserUpdated"
}
