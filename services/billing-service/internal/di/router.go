package di

import (
	"context"
	"encoding/json"
	"net/http"

	"billing-service/internal/api"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/strider2038/pkg/persistence/pgx"
)

func NewAPIRouter(c *Container) http.Handler {
	apiService := api.NewBillingApiService(c.accountRepository, c.operationRepository, c.txManager)
	apiController := api.NewBillingApiController(apiService)

	apiRouter := api.NewRouter(apiController)
	metrics := api.NewMetrics("billing_service")
	apiRouter.Use(func(handler http.Handler) http.Handler {
		return api.MetricsMiddleware(handler, metrics)
	})

	router := NewRouter(c.connection, c.config)
	router.PathPrefix("/api").Handler(apiRouter)
	router.Handle("/metrics", promhttp.Handler())

	return router
}

func NewRouter(connection pgx.Connection, config Config) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"status":"ok"}`))
	})

	router.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		err := connection.Ping(context.Background())
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(`{"status":"ok"}`))
		} else {
			writer.WriteHeader(http.StatusServiceUnavailable)
			writer.Write([]byte(`{"status":"not available"}`))
		}
	})

	router.HandleFunc("/version", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(struct {
			ApplicationName    string `json:"application_name"`
			Environment        string `json:"environment"`
			ApplicationVersion string `json:"application_version"`
		}{
			ApplicationName:    "BillingService",
			Environment:        config.Environment,
			ApplicationVersion: config.Version,
		})
		writer.Write(response)
	})

	return router
}
