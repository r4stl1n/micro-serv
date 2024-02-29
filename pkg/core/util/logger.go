package util

import (
	"go.uber.org/zap"
	"os"
)

type Logger struct{}

// Init create a new zap logger
func (l *Logger) Init() *Logger {

	logger := zap.Must(zap.NewProduction())

	// Check if we have a log level debug
	if os.Getenv("APP_ENV") == "DEV" {
		logger = zap.Must(zap.NewDevelopment())
	}

	zap.ReplaceGlobals(logger)

	return l
}
