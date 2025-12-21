package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"PORT"`
}

func Load() *Config {
	// Production: читаем из environment variables
	viper.AutomaticEnv()

	// Development: пытаемся прочитать .env (если есть)
	viper.SetConfigFile(".env")	// путь к .env файлу
	viper.ReadInConfig()	// чтение переменных окружения

	// Fallback на 8080 если порт не задан
	// viper.SetDefault("PORT", "8080")

	cfg := &Config{}
	viper.Unmarshal(cfg)	// парсинг в структуру

	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	return cfg
}