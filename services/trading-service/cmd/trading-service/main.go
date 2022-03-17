package main

import (
	"context"
	"fmt"
	"log"

	"trading-service/internal/di"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/strider2038/httpserver"
	"github.com/strider2038/ossync"
	"github.com/strider2038/pkg/persistence/pgx"
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

	pool, err := pgxpool.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	connection := pgx.NewPool(pool)
	container, err := di.NewContainer(connection, config)
	if err != nil {
		log.Fatal("failed to create di container: ", err)
	}

	router := di.NewAPIRouter(container)
	address := fmt.Sprintf(":%d", config.Port)
	log.Println("starting application server at", address)

	group := ossync.NewGroup(context.Background())
	group.Go(httpserver.New(address, router).ListenAndServe)
	group.Go(di.NewBillingConsumer(container).Run)
	err = group.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
