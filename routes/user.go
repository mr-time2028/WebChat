package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mr-time2028/WebChat/handlers"
	"net/http"
)

func UserRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/register", handlers.HandlerRepo.Register)
	mux.Post("/login", handlers.HandlerRepo.Login)

	return mux
}
