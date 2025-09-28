package routes

import (
	"net/http"

	"github.com/Romasmi/go-rest-api-template/internal/config"
	"github.com/Romasmi/go-rest-api-template/internal/handlers"
	"github.com/Romasmi/go-rest-api-template/internal/repository"
	"github.com/Romasmi/go-rest-api-template/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUsersRoutes(r *mux.Router, db *pgxpool.Pool, config *config.Config) {
	h := handlers.NewUserHandler(services.NewUserService(repository.NewUserRepository(db)))

	r.HandleFunc("/auth/register", h.Register).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/users", h.ListUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", h.DeleteUser).Methods(http.MethodDelete)
}
