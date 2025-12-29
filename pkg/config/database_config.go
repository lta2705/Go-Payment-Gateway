package config

import (
	"github.com/joho/godotenv"
	"github.com/lta2705/Go-Payment-Gateway/utils"
	"log"
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
		DBHost:            utils.String("DB_HOST", "localhost"),
		DBPort:            utils.String("DB_PORT", "5432"),
		DBUser:            utils.String("DB_USER", "postgres"),
		DBPassword:        utils.String("DB_PASSWORD", "postgres"),
		DBName:            utils.String("DB_NAME", "payment_db"),
		DBSSLMode:         utils.String("DB_SSLMODE", "disable"),
		DBMaxConns:        utils.Int("DB_MAX_OPEN_CONNS", 10),
		DBIdleConn:        utils.Int("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: time.Duration(utils.Int("DB_CONN_MAX_LIFETIME", 300)) * time.Second,
	}
}
