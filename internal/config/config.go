package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string

	DatabaseURL string
	RedisURL    string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	JWTSecret       string
	JWTExpiredHours time.Duration

	RateLimitGlobalCapacity int
	RateLimitGlobalRefill   float64
	RateLimitAuthCapacity   int
	RateLimitAuthRefill     float64
	RateLimitWriteCapacity  int
	RateLimitWriteRefill    float64
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	expHours, _ := strconv.Atoi(getEnv("JWT_EXPIRED_HOURS", "24"))

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", ""),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "taskdb"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		JWTSecret:       getEnv("JWT_SECRET", "default-secret"),
		JWTExpiredHours: time.Duration(expHours) * time.Hour,

		RateLimitGlobalCapacity: getEnvInt("RATE_LIMIT_GLOBAL_CAPACITY", 60),
		RateLimitGlobalRefill:   getEnvFloat("RATE_LIMIT_GLOBAL_REFILL", 1.0),
		RateLimitAuthCapacity:   getEnvInt("RATE_LIMIT_AUTH_CAPACITY", 5),
		RateLimitAuthRefill:     getEnvFloat("RATE_LIMIT_AUTH_REFILL", 5.0/60.0),
		RateLimitWriteCapacity:  getEnvInt("RATE_LIMIT_WRITE_CAPACITY", 20),
		RateLimitWriteRefill:    getEnvFloat("RATE_LIMIT_WRITE_REFILL", 20.0/60.0),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
