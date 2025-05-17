package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"drone/generated"
	"drone/internal/api"
	"drone/internal/config"
	"drone/internal/repository"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Initialize repository
	repo := repository.NewRepository(dbPool)

	// Initialize API handler
	handler := api.NewHandler(repo)

	// Set up Echo server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register API routes
	generated.RegisterHandlers(e, handler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
} 