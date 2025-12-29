package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/internal/healthcheck/handlers"
	"github.com/zeusito/toci/pkg/config"
	"github.com/zeusito/toci/pkg/db"
	"github.com/zeusito/toci/pkg/logger"
	"github.com/zeusito/toci/pkg/router"
)

func main() {
	// Parse flags
	cfgPath := flag.String("config", "resources/config.toml", "Path to the configuration file")
	flag.Parse()

	// Setup logger
	logger.MustConfigure()

	// Load config
	myConfig, err := config.LoadConfigurations(*cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading configurations")
	}

	// Init DB
	myDB := db.MustCreatePooledConnection(myConfig.Database)

	// Init router
	myRouter := router.NewHTTPRouter(myConfig.Server)

	// Health Controller
	_ = handlers.NewHealthController(myRouter.Mux)

	// Start server in background
	go myRouter.Start()

	// Graceful shutdown
	gracefulShutdown(myRouter, myDB)
}

func gracefulShutdown(myRouter *router.HTTPRouter, myDB *db.DatabaseConnection) {
	// Wait for the interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// Signal acquired, starting to shut down all systems
	log.Warn().Msg("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myDB.Close()
	myRouter.Shutdown(ctx)
}
