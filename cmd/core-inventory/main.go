// Package main provides the entry point for mcp-core-inventory
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/vertikon/mcp-core-inventory/internal/adapters/nats"
	"github.com/vertikon/mcp-core-inventory/internal/adapters/postgres"
	"github.com/vertikon/mcp-core-inventory/internal/adapters/redis"
	"github.com/vertikon/mcp-core-inventory/internal/app"
	"github.com/vertikon/mcp-core-inventory/internal/interfaces/http"
	"github.com/vertikon/mcp-core-inventory/internal/observability"
	"github.com/vertikon/mcp-core-inventory/internal/services"
	"github.com/vertikon/mcp-core-inventory/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.Init("info", true); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting mcp-core-inventory")

	// Load configuration from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost/inventory?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	defer rdb.Close()

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		logger.Fatal("Failed to get JetStream context", zap.Error(err))
	}

	// Initialize repositories and services
	ledgerRepo := postgres.NewLedgerRepository(db)
	stockCache := redis.NewStockCache(rdb)
	reservationLock := redis.NewReservationLock(rdb)
	eventPub := nats.NewEventPublisher(js, logger.GetLogger())

	// Initialize use cases
	reserveUseCase := app.NewReserveStockUseCase(
		ledgerRepo,
		reservationLock,
		eventPub,
		logger.GetLogger(),
	)

	confirmUseCase := app.NewConfirmReservationUseCase(
		ledgerRepo,
		ledgerRepo, // ReservationRepository is same as LedgerRepository for now
		eventPub,
		logger.GetLogger(),
	)

	releaseUseCase := app.NewReleaseReservationUseCase(
		ledgerRepo,
		ledgerRepo,
		logger.GetLogger(),
	)

	adjustUseCase := app.NewAdjustStockUseCase(
		ledgerRepo,
		ledgerRepo, // MovementRepository is same as LedgerRepository for now
		logger.GetLogger(),
	)

	queryUseCase := app.NewQueryAvailableUseCase(
		ledgerRepo,
		stockCache,
	)

	// Initialize metrics
	inventoryMetrics := observability.NewInventoryMetrics()

	// Initialize reservation cleanup service
	cleanupService := services.NewReservationCleanupService(
		ledgerRepo,
		ledgerRepo,
		logger.GetLogger(),
		inventoryMetrics,
	)

	// Start cleanup service in background
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	go cleanupService.Start(cleanupCtx)

	// Setup HTTP router
	router := http.NewRouter(
		reserveUseCase,
		confirmUseCase,
		releaseUseCase,
		adjustUseCase,
		queryUseCase,
		logger.GetLogger(),
		inventoryMetrics,
	)

	mux := http.NewServeMux()
	router.SetupRoutes(mux)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("HTTP server starting", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down server")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}

