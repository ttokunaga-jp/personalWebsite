package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
}

var Module = fx.Module("config",
	fx.Provide(load),
)

func load() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "release")

	if err := v.ReadInConfig(); err != nil {
		// Ignore missing file; rely on defaults/env overrides.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
