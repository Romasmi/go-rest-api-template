package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Romasmi/go-rest-api-template/internal/config"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/Romasmi/go-rest-api-template/internal/database"
	"github.com/Romasmi/go-rest-api-template/internal/handlers"
	authMiddleware "github.com/Romasmi/go-rest-api-template/internal/middleware"
	"github.com/Romasmi/go-rest-api-template/internal/repository"
	"github.com/Romasmi/go-rest-api-template/internal/services"
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

	envConfig, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatalf("error while loading config %v", err)
	}

	dbConnection := &database.DbConnection{}
	if err := dbConnection.Connect(envConfig); err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConnection.Close()

	if err := database.RunMigrations("up"); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConnection.DB)

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
	r := mux.NewRouter()

	// Middlewares
	r.Use(ghandlers.RecoveryHandler())
	r.Use(ghandlers.ProxyHeaders)
	r.Use(ghandlers.CORS(
		ghandlers.AllowedOrigins([]string{"*"}),
		ghandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		ghandlers.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
		ghandlers.ExposedHeaders([]string{"Link"}),
		ghandlers.AllowCredentials(),
		ghandlers.MaxAge(300),
	))

	// Public routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go REST API Template!"))
	}).Methods("GET")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1
	api := r.PathPrefix("/api/v1").Subrouter()
	userHandler.RegisterHandlers(api)

	// Protected subrouter
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.Authenticator)
	protected.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected endpoint"))
	}).Methods("GET")

	// Wrap with logging and timeout
	var handler http.Handler = r
	handler = ghandlers.CombinedLoggingHandler(os.Stdout, handler)
	handler = http.TimeoutHandler(handler, 60*time.Second, "Request timed out")

	return handler
}
