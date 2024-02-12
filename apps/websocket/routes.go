package websocket

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", home)
	mux.Get("/ws", wsEndpoint)

	return mux
}
