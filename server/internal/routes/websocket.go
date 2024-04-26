package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mr-time2028/WebChat/internal/handlers"
	"net/http"
)

func (r *RouteRepository) WsRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.HandlerRepo.Home)
	mux.Post("/create_room", handlers.HandlerRepo.CreateRoom)
	mux.Get("/join_room", handlers.HandlerRepo.JoinRoom)
	mux.Get("/get_rooms", handlers.HandlerRepo.GetRooms)
	mux.Post("/get_room_clients", handlers.HandlerRepo.GetRoomClients)

	return mux
}
