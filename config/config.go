package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	dir    = "tax-calculator"
	prefix = "TAX_CALCULATOR"
)

// Config holds all configurations
type Config struct {
	App App
}

// App represents application-specific configurations
type App struct {
	LogPath string `mapstructure:"logPath"`
	Port    int    `mapstructure:"port"`
	Debug   bool   `mapstructure:"debug"`
}

// GetConfig initializes and returns the config
func GetConfig() *Config {
	v := viper.New()
	// Initialize logger
	logger := logrus.New()

	// read flags and files
	v = bindEnv(v)
	v.AutomaticEnv()

	v = bindFlag(v)
	v = bindFile(v)
	v = bindDefault(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		var cfgFileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &cfgFileNotFoundErr) {
			// fallback to default and env vars
			v.AutomaticEnv()
			v = bindDefault(v)
			v = bindEnv(v)
		} else {
			logger.Fatal(err, "error while reading the config file")
		}
	}

	c := &Config{}

	// Unmarshal config file into Config struct
	if err := v.Unmarshal(c); err != nil {
		logger.WithError(err).Fatal("Error while unmarshalling config")
	}

	// Optionally, configure logging path from the config if debug is enabled
	if c.App.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Setting log output
	if c.App.LogPath != "" {
		// Ensures that log file exists
		if _, err := os.Stat(c.App.LogPath); os.IsNotExist(err) {
			err := os.MkdirAll(c.App.LogPath, os.ModePerm)
			if err != nil {
				logger.WithError(err).Fatal("Failed to create log directory")
			}
		}
		file, err := os.OpenFile(filepath.Join(c.App.LogPath, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.WithError(err).Fatal("Failed to open log file")
		}
		logger.SetOutput(file)
	} else {
		// Default to stdout if no log path is set
		logger.SetOutput(os.Stdout)
	}

	return c
}

func bindFlag(v *viper.Viper) *viper.Viper {
	flag := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flag.IntP("App.Port", "p", 8080, "port of the app")
	flag.DurationP("App.RunEvery", "f", 5*time.Minute, "frequency of exporting data")
	flag.StringP("CONFIG_PATH", "c", fmt.Sprintf("/home/adbroker/c_delivery/%s.toml", dir), "location of the config file")

	if err := flag.Parse(os.Args[1:]); err != nil {
		logrus.Fatal("Unexpected error while parsing flags: ", err)
	}

	if err := v.BindPFlags(flag); err != nil {
		logrus.Fatal("Unexpected error while binding flags: ", err)
	}

	return v
}

func bindDefault(v *viper.Viper) *viper.Viper {
	// App defaults
	v.SetDefault("App.Port", 8080)
	v.SetDefault("App.Debug", false)
	v.SetDefault("App.LogPath", "/tmp/tax-calculator")
	return v
}

//nolint:errcheck
func bindEnv(v *viper.Viper) *viper.Viper {
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// CONFIG_PATH is not part of the "Config" struct,
	// so we do not want it prepended
	v.BindEnv("CONFIG_PATH", "CONFIG_PATH")

	// App environment variables
	v.BindEnv("App.Port", "TAX_CALCULATOR_APP_PORT")
	v.BindEnv("App.Debug", "TAX_CALCULATOR_APP_DEBUG")
	v.BindEnv("App.LogPath", "TAX_CALCULATOR_APP_LOG_PATH")
	return v
}

func bindFile(v *viper.Viper) *viper.Viper {
	filePath := v.GetString("CONFIG_PATH")

	if _, err := os.Stat(filePath); err == nil {
		ftype := filepath.Ext(filePath)
		if len(ftype) > 1 {
			ftype = ftype[1:]
		}
		fname := filepath.Base(filePath)
		fname = fname[0 : len(fname)-(len(ftype)+1)]
		fpath := filepath.Dir(filePath)
		v.SetConfigName(fname)
		v.SetConfigType(ftype)
		v.AddConfigPath(fpath)
	} else {
		v.SetConfigName(dir)
		v.SetConfigName("config")
		v.SetConfigType("toml")
		v.AddConfigPath(fmt.Sprintf("/etc/%s/conf.d/", dir))
	}

	return v
}
