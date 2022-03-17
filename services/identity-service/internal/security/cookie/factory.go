package cookie

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Factory struct {
	domain   string
	secure   bool
	lifeTime time.Duration
}

func NewFactory(lifeTime time.Duration, domainURL string) (*Factory, error) {
	u, err := url.Parse(domainURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse domain url: %w", err)
	}

	return &Factory{
		lifeTime: lifeTime,
		domain:   u.Host,
		secure:   u.Scheme == "https",
	}, nil
}

func (factory *Factory) Create(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   factory.domain,
		Expires:  time.Now().Add(factory.lifeTime),
		Secure:   factory.secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
