package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Romasmi/go-rest-api-template/internal/application"
)

func main() {
	app := &application.App{}
	err := app.InitApp(".")
	if err != nil {
		fmt.Printf("error while app initialization: %v", err)
		return
	}
	defer app.OnStop()

	app.Run()
}

func setupRouter(logger *log.Logger) http.Handler {
	// This is a temporary implementation that will be replaced
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Go REST API Template!")
	})

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	return mux
}
