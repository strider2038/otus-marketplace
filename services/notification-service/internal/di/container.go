package di

import (
	"os"

	"notification-service/internal/monitoring"
	"notification-service/internal/notifications"
	"notification-service/internal/postgres"
	"notification-service/internal/postgres/database"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Container struct {
	config Config

	dbConnection *pgxpool.Pool
	db           *database.Queries
	logger       zerolog.Logger
	metrics      *monitoring.Metrics

	notifications notifications.Repository
	users         notifications.UserRepository
}

func NewContainer(connection *pgxpool.Pool, config Config) *Container {
	c := &Container{dbConnection: connection, config: config}
	c.db = database.New(connection)
	c.logger = zerolog.New(os.Stdout)
	c.metrics = monitoring.NewMetrics("notification_service")
	c.notifications = postgres.NewNotificationRepository(c.db)
	c.users = postgres.NewUserRepository(c.db)

	return c
}
