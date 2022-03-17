package di

type Config struct {
	DatabaseURL      string `env:"DATABASE_URL" env-required:"true"`
	KafkaConsumerURL string `env:"KAFKA_CONSUMER_URL" env-required:"true"`
	KafkaProducerURL string `env:"KAFKA_PRODUCER_URL" env-required:"true"`

	Environment string `env:"APP_ENV" env-default:"unknown"`
	Version     string
	Port        int `env:"PORT" env-default:"8000"`
}
