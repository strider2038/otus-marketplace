package jwt

import (
	"fmt"
	"time"

	"identity-service/internal/users"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

type Issuer struct {
	privateKey []byte
	lifeTime   time.Duration
}

func NewIssuer(privateKey []byte, lifeTime time.Duration) *Issuer {
	return &Issuer{privateKey: privateKey, lifeTime: lifeTime}
}

func (issuer *Issuer) Issue(user *users.User) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(issuer.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse RSA key: %w", err)
	}

	tokenID := uuid.Must(uuid.NewV4())
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(issuer.lifeTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "identity-service",
		},
		UserID: user.ID,
		Email:  user.Email,
	})
	token.Header["kid"] = KeyID

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
