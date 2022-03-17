package jwt

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

const KeyID = "auth"

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}
