package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/yourusername/go-rest-api-template/internal/database"
	"github.com/yourusername/go-rest-api-template/internal/handlers"
	authMiddleware "github.com/yourusername/go-rest-api-template/internal/middleware"
	"github.com/yourusername/go-rest-api-template/internal/repository"
	"github.com/yourusername/go-rest-api-template/internal/services"
)

// @title Go REST API Template
// @version 1.0
// @description A RESTful API template using Go
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()

	logger := log.New(os.Stdout, "API: ", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authMiddleware.InitAuth()

	if err := database.Connect(); err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations("up"); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(database.DB)

	userService := services.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)

	router := setupRouter(logger, userHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on port %s: %v\n", port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server gracefully stopped")
}

func setupRouter(logger *log.Logger, userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go REST API Template!"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterHandlers(r)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(authMiddleware.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("This is a protected endpoint"))
			})
		})
	})

	return r
}
