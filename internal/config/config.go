package config

import "os"

type Config struct {
	GRPCPort  string
	JWTSecret string
	DBConnStr string
	RedisUrl  string
}

func Load() *Config {
	return &Config{
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		JWTSecret: getEnv("JWT_SECRET", "secret_password"),
		DBConnStr: getEnv("DB_CONN", "postgres://user:pass@pgsql:5432/auth?sslmode=disable"),
		RedisUrl:  getEnv("REDIS_URL", "redisL:6379"),
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
