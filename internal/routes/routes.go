package routes

import (
	"encoding/json"

	"github.com/Romasmi/go-rest-api-template/internal/config"
	authMiddleware "github.com/Romasmi/go-rest-api-template/internal/middleware"
	ghandlers "github.com/gorilla/handlers"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotFoundResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(r *mux.Router, db *pgxpool.Pool, config *config.Config) {
	if r == nil {
		panic("r must be initialized before routes registration")
	}

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Use(ghandlers.RecoveryHandler())
	r.Use(ghandlers.ProxyHeaders)
	r.Use(ghandlers.CORS(
		ghandlers.AllowedOrigins([]string{"*"}),
		ghandlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions}),
		ghandlers.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
		ghandlers.ExposedHeaders([]string{"Link"}),
		ghandlers.AllowCredentials(),
		ghandlers.MaxAge(300),
	))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go REST API Template!"))
	}).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	RegisterUsersRoutes(api, db, config)

	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.Authenticator)
	protected.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected endpoint"))
	}).Methods(http.MethodGet)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	response := &NotFoundResponse{Error: "Route Not found"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
