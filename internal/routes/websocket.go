package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mr-time2028/WebChat/internal/handlers"
	"net/http"
)

func WsRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.HandlerRepo.Home)
	mux.Get("/ws", handlers.HandlerRepo.WsEndpoint)

	return mux
}
