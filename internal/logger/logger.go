package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Declare a global logger
var Log zerolog.Logger

// Initialize the logger globally
func InitLogger() {
	logLevel := zerolog.InfoLevel // Default log level
	if viper.GetBool("App.Debug") {
		logLevel = zerolog.DebugLevel
	}

	// Create a logger instance
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)

	// Optionally, if you want to write to a file
	logPath := viper.GetString("App.LogPath")
	if logPath != "" {
		if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
			Log.Fatal().Err(err).Msg("Failed to create log directory")
		}
		file, err := os.OpenFile(logPath+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			Log.Fatal().Err(err).Msg("Failed to open log file")
		}
		Log = zerolog.New(file).With().Timestamp().Logger().Level(logLevel)
	}
}
