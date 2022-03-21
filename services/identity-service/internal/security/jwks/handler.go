package jwks

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"identity-service/internal/security/jwt"

	"github.com/lestrrat-go/jwx/jwk"
)

type Handler struct {
	keys []jwk.Key
}

func NewHandler(publicKey []byte) (*Handler, error) {
	publicPEM, _ := pem.Decode(publicKey)
	if publicPEM == nil {
		return nil, errors.New("failed to parse RSA public key")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(publicPEM.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA public key: %w", err)
	}

	key, err := jwk.New(parsedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create symmetric key: %w", err)
	}
	key.Set(jwk.KeyTypeKey, "RSA")
	key.Set(jwk.AlgorithmKey, "RSA512")
	key.Set(jwk.KeyIDKey, jwt.KeyID)

	return &Handler{keys: []jwk.Key{key}}, nil
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(struct {
		Keys []jwk.Key `json:"keys"`
	}{
		Keys: h.keys,
	})
	writer.Write(bytes)
}
