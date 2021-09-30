package util

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func RunCleanerScheduler(logger *zap.Logger) {
	if viper.GetBool("config.cleanup.enable") {
		logger.Info("setup scheduler")
		scheduler := gocron.NewScheduler(time.Local)
		scheduler.Every(viper.GetString("config.cleanup.interval")).Do(func() {
			cleanTempDirectory(logger)
		})
		scheduler.StartAsync()
	} else {
		logger.Info("scheduler not enabled")
	}
}

func cleanTempDirectory(logger *zap.Logger) {
	logger.Info("start clean up the temp directory")
	logger.Debug("clean up " + viper.GetString("config.logstash.pipeline_dir"))

	err := os.RemoveAll(viper.GetString("config.logstash.pipeline_dir"))
	if err != nil {
		logger.Error("an error occured",
			zap.String("stacktrace", err.Error()))
		os.Exit(1)
	}

	logger.Info("cleaned up successfully")
}
