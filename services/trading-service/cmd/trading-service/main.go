package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"trading-service/internal/di"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/strider2038/httpserver"
	"github.com/strider2038/ossync"
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
	config.LockTimeout = 5 * time.Second

	container, err := di.NewContainer(config)
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
