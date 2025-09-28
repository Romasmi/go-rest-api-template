package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Romasmi/go-rest-api-template/internal/config"
	"github.com/Romasmi/go-rest-api-template/internal/database"
	authMiddleware "github.com/Romasmi/go-rest-api-template/internal/middleware"
	"github.com/Romasmi/go-rest-api-template/internal/routes"
	ghandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

type App struct {
	config *config.Config
	dbConn *database.DbConnection
	router *mux.Router
	logger *log.Logger
}

func (app *App) InitApp(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading config: %v\n", err)
	}
	app.config = envConfig

	dbConn := &database.DbConnection{Config: envConfig}
	err = dbConn.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}
	app.dbConn = dbConn
	app.router = mux.NewRouter()
	routes.RegisterRoutes(app.router, app.dbConn.DB, app.config)

	authMiddleware.InitAuth(envConfig)

	app.logger = log.New(os.Stdout, "API: ", log.LstdFlags)

	if err := database.RunMigrations("up"); err != nil {
		app.logger.Fatalf("Failed to run migrations: %v", err)
	}

	return nil
}

func (app *App) OnStop() {
	app.dbConn.Close()
}

func (app *App) Run() {
	var handler http.Handler = app.router
	handler = ghandlers.CombinedLoggingHandler(os.Stdout, handler)
	handler = http.TimeoutHandler(handler, 60*time.Second, "Request timed out")

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(int(app.config.Server.Port)),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	app.logger.Println("Server gracefully stopped")
}
