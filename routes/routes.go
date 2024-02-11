package routes

import (
	"github.com/mr-time2028/WebChat/apps/user"
	"github.com/mr-time2028/WebChat/apps/websocket"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Mount("/", websocket.Routes())
	mux.Mount("/users", user.Routes())

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/web/static/*", http.StripPrefix("/web/static", fileServer))

	return mux
}
