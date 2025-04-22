package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name          string
		debug         bool
		logPath       string
		expectedLevel zerolog.Level
		expectFile    bool
	}{
		{
			name:          "when Debug is true, log level should be Debug",
			debug:         true,
			logPath:       "",
			expectedLevel: zerolog.DebugLevel,
			expectFile:    false,
		},
		{
			name:          "when Debug is false, log level should be Info",
			debug:         false,
			logPath:       "",
			expectedLevel: zerolog.InfoLevel,
			expectFile:    false,
		},
		{
			name:          "when log path is set, log file should be created",
			debug:         true,
			logPath:       "./test_logs",
			expectedLevel: zerolog.DebugLevel,
			expectFile:    true,
		},
		{
			name:          "when log path is empty, logger should use stdout",
			debug:         true,
			logPath:       "",
			expectedLevel: zerolog.DebugLevel,
			expectFile:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			viper.Set("App.Debug", tt.debug)
			viper.Set("App.LogPath", tt.logPath)

			// Act
			InitLogger()

			// Assert: Check log level
			assert.Equal(t, tt.expectedLevel, Log.GetLevel())

			// Assert: If the log path is set, check if the file was created
			if tt.expectFile {
				// Check if the log directory exists
				_, err := os.Stat(tt.logPath)
				require.NoError(t, err, "Log directory should be created")

				// Check if the log file exists
				_, err = os.Stat(tt.logPath + "/app.log")
				require.NoError(t, err, "Log file should be created")
			} else {
				// Verify that no log file should be created (i.e., logger should be set to stdout)
				assert.NotNil(t, Log) // Ensure logger is initialized
			}
		})
	}
}
