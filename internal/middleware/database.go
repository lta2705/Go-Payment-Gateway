package middleware

import (
	"github.com/lta2705/Go-Payment-Gateway/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func SetupDatabase(cfg *config.DBConfig) *gorm.DB {
	logger := SetupLogger()
	//Create connection
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=" + cfg.DBSSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.DBMaxConns)
	sqlDB.SetMaxIdleConns(cfg.DBIdleConn)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	logger.Info("Database connection established")
	return db
}
