package config

import (
	"os"
	"strings"
)

type JWTConfig struct {
	Secret string
	Issuer string
}

type Config struct {
	Addr string
	JWT  JWTConfig
}

func valueOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	return Config{
		Addr: valueOrDefault("UBOT_HTTP_ADDR", ":8080"),
		JWT: JWTConfig{
			Secret: valueOrDefault("UBOT_JWT_SECRET", "Ubot-can-be-better"),
			Issuer: valueOrDefault("UBOT_JWT_ISSUER", "Ubot"),
		},
	}
}
