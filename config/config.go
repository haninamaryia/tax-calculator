package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	dir    = "tax-calculator"
	prefix = "TAX_CALCULATOR"
)

// Config hold All configurations
type Config struct {
	App App
}

// App represents application specific configurations
type App struct {
	LogPath string `mapstructure:"logPath"`
	Port    int    `mapstructure:"port"`
	Debug   bool   `mapstructure:"debug"`
}

// GetConfig initializes and return the config
func GetConfig() *Config {
	v := viper.New()
	// read flags and files
	v = bindEnv(v)
	v.AutomaticEnv()

	v = bindFlag(v)
	v = bindFile(v)
	v = bindDefault(v)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("error while reading the config file %w", err)
	}
	c := &Config{}

	if err := v.Unmarshal(c); err != nil {
		log.Fatal("error while unmarshalling config", err)
	}

	return c
}

func bindFlag(v *viper.Viper) *viper.Viper {
	flag := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flag.IntP("App.Port", "p", 20000, "port of the app")
	flag.DurationP("App.RunEvery", "f", 5*time.Minute, "frequency of exporting data")
	flag.StringP("CONFIG_PATH", "c", fmt.Sprintf("/home/adbroker/c_delivery/%s.toml", dir), "location of the config file")

	if err := flag.Parse(os.Args[1:]); err != nil {
		log.Fatal("unexpected error while parsing flags", err)
	}

	if err := v.BindPFlags(flag); err != nil {
		log.Fatal("unexpected error while binding flags", err)
	}

	return v
}

func bindDefault(v *viper.Viper) *viper.Viper {
	// app
	v.SetDefault("App.Port", 20000)
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

	// app
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
