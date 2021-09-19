package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/misbahulard/logstash-pipeline-tester/config"
	"github.com/misbahulard/logstash-pipeline-tester/handler"
	"github.com/misbahulard/logstash-pipeline-tester/logging"
	mdl "github.com/misbahulard/logstash-pipeline-tester/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()

	logger := logging.InitLogger()
	defer logger.Sync()

	logger.Info("Starting the application",
		zap.String("address", viper.GetString("config.address")),
		zap.String("port", viper.GetString("config.port")))

	serve := viper.GetString("config.address") + ":" + viper.GetString("config.port")

	r := chi.NewRouter()
	r.Use(mdl.LoggerMiddleware(logger))
	r.Use(cors.AllowAll().Handler)

	r.Get("/", handler.RootHandler(logger))
	r.Get("/health", handler.HealthCheckHandler(logger))
	r.Post("/api/pipeline", handler.PipelineHandler(logger))

	http.ListenAndServe(serve, r)
}
