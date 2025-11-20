package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type DBConfig struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	ServerPort        string
	DBSSLMode         string
	DBMaxConns        int
	DBIdleConn        int
	DBIdleTimeout     time.Duration
	DBConnMaxLifetime time.Duration
}

func LoadDBConfig() *DBConfig {
	// Load .env file
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return &DBConfig{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "payment_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxConns:        getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
		DBIdleConn:        getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 300)) * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
