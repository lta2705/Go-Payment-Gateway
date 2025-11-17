package main

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/app"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"go.uber.org/zap"
	"os"
)

func main() {
	// Set up logger
	logger := middleware.CreateLogger()
	defer logger.Sync()

	// Initialize the app
	app, err := app.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	logger.Info("Starting server", zap.String("port", port))
	if err := app.Run(":" + port); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
		panic("Failed to run server: " + err.Error())
	}

}
