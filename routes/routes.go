package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Mount("/", WsRoutes())
	mux.Mount("/users", UserRoutes())

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/web/static/*", http.StripPrefix("/web/static", fileServer))

	return mux
}
