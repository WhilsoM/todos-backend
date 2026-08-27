package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret          string `env:"JWT_SECRET" env-required:"true"`
	DatabaseURL        string `env:"DATABASE_URL" env-required:"true"`
	Port               string `env:"PORT" env-default:":8080"`
	RedisPort          string `env:"REDIS_PORT" env-default:"localhost:6379"`
	UserEventKafkaPort string `env:"USER_EVENT_PRODUCER_KAFKA_PORT" env-default:"localhost:9092"`
}

func MustConfig() *Config {
	_ = godotenv.Load()

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
