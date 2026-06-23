package main

import (
	"log"
	nethttp "net/http"

	"github.com/askii12/launchpad/internal/config"
	"github.com/askii12/launchpad/internal/database"
	httptransport "github.com/askii12/launchpad/internal/transport/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load application configuration from environment variables.
	// This allows us to change behavior without rebuilding the app.
	cfg := config.Load()

	// Initialize PostgreSQL connection using GORM.
	// This creates a connection pool used across the entire application.
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	// Extract underlying sql.DB to perform low-level health checks.
	// GORM abstracts this, but we still need direct access for Ping.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB instance:", err)
	}

	// Verify database connectivity at startup.
	// Fail fast if DB is unreachable to avoid running a broken service.
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("database connection check failed:", err)
	}

	log.Println("database connection established")

	// Initialize HTTP router (Gin engine).
	// This will act as the main request handler.
	router := gin.Default()

	// Register all HTTP routes in a separate transport layer.
	// Keeps main.go clean and enforces separation of concerns.
	httptransport.RegisterRoutes(router, db)

	// Configure HTTP server explicitly instead of using router.Run().
	// This is required for future graceful shutdown implementation.
	server := &nethttp.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	log.Println("http server starting on port", cfg.AppPort)

	// Start HTTP server (blocking call).
	// In production this would be wrapped with graceful shutdown logic.
	if err := server.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
		log.Fatal("server failed:", err)
	}
}
