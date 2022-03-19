package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"identity-service/internal/di"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/strider2038/httpserver"
)

var (
	version = ""
)

func main() {
	config := di.Config{Version: version}
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		log.Fatal("invalid config:", err)
	}
	dbConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		log.Fatal("failed to parse postgres config:", err)
	}

	logger := zerolog.New(os.Stdout)
	dbConfig.ConnConfig.Logger = zerologadapter.NewLogger(logger)
	dbConfig.ConnConfig.LogLevel = pgx.LogLevelInfo

	connection, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	router, err := di.NewRouter(connection, config)
	if err != nil {
		log.Fatal("failed to create router: ", err)
	}

	address := fmt.Sprintf(":%d", config.Port)
	log.Println("starting HTTP server at", address)
	server := httpserver.New(address, router)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	err = server.ListenAndServe(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
