package di

import "time"

type Config struct {
	DatabaseURL     string        `env:"DATABASE_URL" env-required:"true"`
	KafkaURL        string        `env:"KAFKA_URL" env-required:"true"`
	PrivateKey      []byte        `env:"PRIVATE_KEY" env-required:"true"`
	PublicKey       []byte        `env:"PUBLIC_KEY" env-required:"true"`
	DomainURL       string        `env:"DOMAIN_URL" env-required:"true"`
	SessionLifeTime time.Duration `env:"SESSION_LIFE_TIME" env-default:"24h"`

	Port        int    `env:"PORT" env-default:"8000"`
	Environment string `env:"APP_ENV" env-default:"unknown"`
	Version     string
}
