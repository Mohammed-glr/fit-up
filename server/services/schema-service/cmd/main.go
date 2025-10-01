package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/handlers"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/service"
)

func main() {
	log.Println("=================================================")
	log.Println("FitUp Schema Service Starting...")
	log.Println("=================================================")

	// Load configuration from environment
	config := loadConfig()

	// Initialize database connection
	db, err := initDatabase(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Database connection established")

	// Initialize repository layer
	repo := repository.NewStore(db)
	log.Println("✓ Repository layer initialized")

	// Initialize service layer
	svc := service.NewService(repo)
	log.Println("✓ Service layer initialized")

	// Initialize HTTP router
	r := setupRouter()

	// Register handlers
	registerHandlers(r, svc)
	log.Println("✓ Handlers registered")

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("✓ Schema Service listening on port %s", config.Port)
		log.Println("=================================================")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n=================================================")
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server exited gracefully")
	log.Println("=================================================")
}

// Config holds application configuration
type Config struct {
	Port        string
	DatabaseURL string
	Environment string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	config := &Config{
		Port:        getEnv("PORT", "8083"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	return config
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// initDatabase initializes database connection pool
func initDatabase(dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Configure pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

// setupRouter configures the HTTP router with middleware
func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware - basic configuration
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	return r
}

// registerHandlers registers all API handlers
func registerHandlers(r *chi.Mux, svc service.SchemaService) {
	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"schema-service","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Exercise routes
		exerciseHandler := handlers.NewExerciseHandler(svc.Exercises())
		exerciseHandler.RegisterRoutes(r)

		// Workout routes
		workoutHandler := handlers.NewWorkoutHandler(svc.Workouts())
		workoutHandler.RegisterRoutes(r)

		// Fitness Profile routes
		fitnessProfileHandler := handlers.NewFitnessProfileHandler(svc.FitnessProfiles())
		fitnessProfileHandler.RegisterRoutes(r)

		// Workout Session routes
		workoutSessionHandler := handlers.NewWorkoutSessionHandler(svc.WorkoutSessions())
		workoutSessionHandler.RegisterRoutes(r)

		// Plan Generation routes
		planGenerationHandler := handlers.NewPlanGenerationHandler(svc.PlanGeneration())
		planGenerationHandler.RegisterRoutes(r)

		// Performance Analytics routes
		performanceAnalyticsHandler := handlers.NewPerformanceAnalyticsHandler(svc.PerformanceAnalytics())
		performanceAnalyticsHandler.RegisterRoutes(r)
	})
}
