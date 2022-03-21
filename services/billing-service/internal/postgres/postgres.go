package postgres

import (
	"context"

	"billing-service/internal/postgres/database"

	postgres "github.com/strider2038/pkg/persistence/pgx"
)

func queries(ctx context.Context, conn postgres.Connection) *database.Queries {
	return database.New(conn.Scope(ctx))
}
