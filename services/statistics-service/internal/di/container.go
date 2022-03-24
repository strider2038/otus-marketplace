package di

import (
	"os"

	"statistics-service/internal/monitoring"
	"statistics-service/internal/postgres"
	"statistics-service/internal/postgres/database"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Container struct {
	connection *pgxpool.Pool
	db         *database.Queries
	config     Config
	logger     zerolog.Logger
	metrics    *monitoring.Metrics

	dailyDeals      *postgres.DailyDealsRepository
	totalDailyDeals *postgres.TotalDailyDealsRepository
	top10Deals      *postgres.Top10DealsRepository
}

func NewContainer(connection *pgxpool.Pool, config Config) *Container {
	c := &Container{connection: connection, config: config}
	c.db = database.New(connection)
	c.logger = zerolog.New(os.Stdout)
	c.metrics = monitoring.NewMetrics("statistics_service")
	c.dailyDeals = postgres.NewDailyDealsRepository(c.db)
	c.totalDailyDeals = postgres.NewTotalDailyDealsRepository(c.db)
	c.top10Deals = postgres.NewTop10DealsRepository(c.db)

	return c
}
