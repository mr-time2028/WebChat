package websocket

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", http.HandlerFunc(Home))
	mux.Get("/ws", http.HandlerFunc(WsEndpoint))

	return mux
}
