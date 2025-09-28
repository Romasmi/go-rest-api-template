package routes

import (
	"encoding/json"
	"github.com/Romasmi/go-rest-api-template/internal/config"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotFoundResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(router *mux.Router, db *pgxpool.Pool, config *config.Config) {
	if router == nil {
		panic("router must be initialized before routes registration")
	}

	RegisterUsersRoutes(router, db, config)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
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
