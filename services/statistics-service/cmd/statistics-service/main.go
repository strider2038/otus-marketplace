package main

import (
	"context"
	"fmt"
	"log"
	"os"

	di2 "statistics-service/internal/di"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/strider2038/httpserver"
	"github.com/strider2038/ossync"
)

var (
	version = ""
)

func main() {
	config := di2.Config{Version: version}
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

	pool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	container := di2.NewContainer(pool, config)

	router := di2.NewAPIRouter(container)
	address := fmt.Sprintf(":%d", config.Port)
	log.Println("starting application server at", address)

	group := ossync.NewGroup(context.Background())
	group.Go(httpserver.New(address, router).ListenAndServe)
	group.Go(di2.NewTradingConsumer(container).Run)
	err = group.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
