package logger

import (
	"log"
	"log/slog"
	"openmyth/messgener/config"
	"openmyth/messgener/pkg/common"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetLoggerGlobal() {
	var slogHandler slog.Handler
	switch config.GetGlobalConfig().Env {
	case common.EnvironmentProduction, common.EnvironmentStaging, common.EnvironmentDev:
		f, err := os.OpenFile(
			config.GetGlobalConfig().Log.FileOutput,
			os.O_RDWR|os.O_CREATE|os.O_APPEND,
			0666,
		)
		if err != nil {
			log.Fatalf("unable to open log file output: %v", err)
		}
		slogHandler = slog.NewJSONHandler(f, &slog.HandlerOptions{
			AddSource: true,
		})
	default:
		slogHandler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			TimeFormat: time.RFC3339Nano,
		})
	}

	logger := slog.New(slogHandler)

	slog.SetDefault(logger)
}
