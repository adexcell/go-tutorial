package config

import "github.com/spf13/viper"

type Config struct {
	Port string `mapstructure:"PORT"`
}

func Load() *Config {
	viper.SetConfigFile(".env")	// путь к .env файлу
	viper.ReadInConfig()	// чтение переменных окружения
	cfg := &Config{}
	viper.Unmarshal(cfg)	// парсинг в структуру
	return cfg
}