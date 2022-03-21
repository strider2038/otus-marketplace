package di

import (
	"context"
	"encoding/json"
	"net/http"

	"history-service/internal/api"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewAPIRouter(c *Container) http.Handler {
	apiService := api.NewHistoryApiService(c.dealRepository)
	apiController := api.NewHistoryApiController(apiService)

	apiRouter := api.NewRouter(apiController)
	metrics := api.NewMetrics("history_service")
	apiRouter.Use(func(handler http.Handler) http.Handler {
		return api.MetricsMiddleware(handler, metrics)
	})

	router := NewRouter(c.connection, c.config)
	router.PathPrefix("/api").Handler(apiRouter)
	router.Handle("/metrics", promhttp.Handler())

	return router
}

func NewRouter(connection *pgxpool.Pool, config Config) *mux.Router {
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
			ApplicationName:    "TradingService",
			Environment:        config.Environment,
			ApplicationVersion: config.Version,
		})
		writer.Write(response)
	})

	return router
}
