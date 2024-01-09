package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mr-time2028/WebChat/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/web/static/*", http.StripPrefix("/web/static", fileServer))

	return mux
}
