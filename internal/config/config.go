package config

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Env        string           `mapstructure:"env" validate:"required"`
	HTTPServer HTTPServerConfig `mapstructure:"http_server" validate:"required"`
	Auth       AuthConfig       `mapstructure:"auth" validate:"required"`
	Gin        GinConfig        `mapstructure:"gin" validate:"required"`
	Postgres   PostgresConfig   `mapstructure:"postgres" validate:"required"`
	Redis      RedisConfig      `mapstructure:"redis" validate:"required"`
	RabbitMQ   RabbitMQConfig   `mapstructure:"rabbitmq" validate:"required"`
	Kafka      KafkaConfig      `mapstructure:"kafka" validate:"required"`
	Logger     LoggerConfig     `mapstructure:"logger"`
}

type HTTPServerConfig struct {
	Addr            string        `mapstructure:"addr" validate:"required"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" validate:"required"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout" validate:"required"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes" validate:"required"`
}

type AuthConfig struct {
	JWTSecret string        `mapstructure:"jwt_secret" validate:"required"`
	TokenTTL  time.Duration `mapstructure:"token_ttl" validate:"required"`
}

type GinConfig struct {
	Mode string `mapstructure:"mode" validate:"required"`
}

type PostgresConfig struct {
	DSN             string        `mapstructure:"dsn" validate:"required"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"required"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"required"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Addr         string        `mapstructure:"addr" validate:"required"`
	Password     string        `mapstructure:"password" validate:"required"`
	DB           int           `mapstructure:"db" validate:"required"`
	MinIdleConns int           `mapstructure:"min_idle_conns" validate:"required"`
	PoolSize     int           `mapstructure:"pool_size" validate:"required"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" validate:"required"`
	TTL          time.Duration `mapstructure:"ttl" validate:"required"`
}

type RabbitMQConfig struct {
	URL           string `mapstructure:"url" validate:"required"`
	Exchange      string `mapstructure:"exchange" validate:"required"`
	Kind          string `mapstructure:"kind" validate:"required"`
	DeliveryMode  int    `mapstructure:"delivery_mode" validate:"required"`
	PrefetchCount int    `mapstructure:"prefetch_count" validate:"required"`
}

type KafkaConfig struct {
	Brokers             []string      `mapstructure:"brokers" validate:"required"`
	ClientID            string        `mapstructure:"client_id" validate:"required"`
	ConsumerGroup       string        `mapstructure:"consumer_group" validate:"required"`
	Topic               string        `mapstructure:"topic" validate:"required"`
	ConsumerWorkerCount int           `mapstructure:"consumer_worker_count" validate:"required"`
	RetryMax            int           `mapstructure:"retry_max" validate:"required"`
	RequiredAcks        int           `mapstructure:"required_acks" validate:"required"`
	MaxWaitTime         time.Duration `mapstructure:"max_wait_time" validate:"required"`
	BatchSize           int           `mapstructure:"batch_size" validate:"required"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level" default:"info"`
	JSONFormat bool   `mapstructure:"json_format" validate:"required"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфига: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига в структуру: %w", err)
	}

	return &cfg, nil
}
