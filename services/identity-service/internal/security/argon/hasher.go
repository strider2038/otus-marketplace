package argon

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

type Hasher struct{}

func (hasher Hasher) Hash(password string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(encoded), nil
}

func (hasher Hasher) Verify(password, hash string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(hash))
}
