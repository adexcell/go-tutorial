package config

import "github.com/spf13/viper"

type Config struct {
	Port string `mapstructure:"PORT"`
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	cfg := &Config{}
	viper.Unmarshal(cfg)
	return cfg
}