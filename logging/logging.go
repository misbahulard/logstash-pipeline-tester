package logging

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	encoderConfig := getEncoder()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writer := getLogWriter()

	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func getEncoder() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return encoderConfig
}

func getLogWriter() zapcore.WriteSyncer {
	if viper.GetBool("config.log.file.enable") {
		writer := &lumberjack.Logger{
			Filename:   viper.GetString("config.log.file.path"),
			MaxSize:    viper.GetInt("config.log.file.max_size"),
			MaxBackups: viper.GetInt("config.log.file.max_backup"),
			MaxAge:     viper.GetInt("config.log.file.max_age"),
			Compress:   viper.GetBool("config.log.file.compress"),
		}

		return zapcore.AddSync(writer)
	}

	writer := zapcore.Lock(os.Stdout)
	return writer
}
