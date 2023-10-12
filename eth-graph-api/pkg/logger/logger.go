package logger

import (
	"go.uber.org/zap"
	"log"
)

// Log is the generalized logger for the application
var Log *zap.SugaredLogger

// Initialize initializes the logger to be used across the application
func Initialize() error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	Log = log.Sugar()

	return nil
}

// Error logs error messages with key-value pairs
func Error(message string, keysAndValues ...interface{}) {
	if Log == nil {
		err := Initialize()
		if err != nil {
			panic(err)
		}
		defer func(Log *zap.SugaredLogger) {
			err := Log.Sync()
			if err != nil {
				log.Printf("Error syncing logger: %v", err)
			}
		}(Log)
	}

	Log.Errorw(message, keysAndValues...)
}
