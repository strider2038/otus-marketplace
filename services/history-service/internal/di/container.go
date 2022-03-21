package di

import (
	"os"

	"history-service/internal/history"
	"history-service/internal/postgres"
	"history-service/internal/postgres/database"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Container struct {
	connection *pgxpool.Pool
	db         *database.Queries
	config     Config
	logger     zerolog.Logger

	dealRepository history.DealRepository
}

func NewContainer(connection *pgxpool.Pool, config Config) *Container {
	c := &Container{connection: connection, config: config}
	c.db = database.New(connection)
	c.logger = zerolog.New(os.Stdout)
	c.dealRepository = postgres.NewDealRepository(c.db)

	return c
}
