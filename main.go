package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/misbahulard/logstash-pipeline-tester/config"
	"github.com/misbahulard/logstash-pipeline-tester/handler"
	"github.com/misbahulard/logstash-pipeline-tester/logging"
	mdl "github.com/misbahulard/logstash-pipeline-tester/middleware"
	"github.com/misbahulard/logstash-pipeline-tester/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//go:embed web/logstash-pipeline-tester/dist/*
var spa embed.FS

func main() {
	fsSub, _ := fs.Sub(spa, "web/logstash-pipeline-tester/dist")

	config.InitConfig()

	logger := logging.InitLogger()
	defer logger.Sync()

	logger.Info("Starting the application",
		zap.String("address", viper.GetString("config.address")),
		zap.String("port", viper.GetString("config.port")))

	util.RunCleanerScheduler(logger)

	serve := viper.GetString("config.address") + ":" + viper.GetString("config.port")

	r := chi.NewRouter()
	r.Use(mdl.LoggerMiddleware(logger))
	r.Use(cors.AllowAll().Handler)

	r.Get("/api/health", handler.HealthCheckHandler(logger))
	r.Post("/api/pipeline", handler.PipelineHandler(logger))

	// you need to build the web first using `yarn build`
	r.Get("/*", func(rw http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(fsSub)).ServeHTTP(rw, r)
	})

	http.ListenAndServe(serve, r)
}

func FileServer(router *chi.Mux, resources embed.FS) {
	fs := http.FileServer(http.FS(resources))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
	})
}
