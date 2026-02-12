package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/vahiiiid/go-rest-api-boilerplate/api/docs"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/config"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/db"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/imoveis"
)

func main() {
	// Parse command-line flags (organization ID is no longer required)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Connect to database
	database, err := db.NewPostgresDBFromDatabaseConfig(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	sqlDB, err := database.DB()
	if err != nil {
		logger.Error("Failed to get database connection", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
	}()

	logger.Info("Connected to database successfully")

	// Initialize services
	imoveisRepo := imoveis.NewRepository(database)
	imoveisService := imoveis.NewService(imoveisRepo)
	// Organization ID is now taken from the external API data
	imoveisImportService := imoveis.NewImportService(imoveisService, &cfg.ExternalAPI)

	logger.Info("Starting import of properties from external API")

	// Run import
	ctx := context.Background()
	if err := imoveisImportService.ImportPublishedProperties(ctx); err != nil {
		logger.Error("Import completed with message", "result", err.Error())
	}

	logger.Info("Import process finished")
}
