package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		unsetEnv    []string
		expectedCfg Config
	}{
		{
			name: "WithDefaults",
			unsetEnv: []string{
				"TAX_CALCULATOR_APP_PORT",
				"TAX_CALCULATOR_APP_DEBUG",
				"TAX_CALCULATOR_APP_LOG_PATH",
				"CONFIG_PATH",
			},
			expectedCfg: Config{App{Port: 8080,
				Debug:   false,
				LogPath: "/tmp/tax-calculator",
			}},
		},
		{
			name: "WithEnvVars",
			envVars: map[string]string{
				"TAX_CALCULATOR_APP_PORT":  "8080",
				"TAX_CALCULATOR_APP_DEBUG": "true",
			},
			expectedCfg: Config{App{Port: 8080,
				Debug:   true,
				LogPath: "/tmp/tax-calculator"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset env vars
			for _, key := range tt.unsetEnv {
				os.Unsetenv(key)
			}
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k) // Clean up
			}

			// Reset CLI args to simulate config path flag
			//os.Args = []string{"cmd", "--CONFIG_PATH=./doesnotexist.toml"}

			cfg := GetConfig()

			assert.NotNil(t, cfg)
			assert.Equal(t, tt.expectedCfg.App.Port, cfg.App.Port)
			assert.Equal(t, tt.expectedCfg.App.Debug, cfg.App.Debug)
			assert.Equal(t, tt.expectedCfg.App.LogPath, cfg.App.LogPath)
		})
	}
}
