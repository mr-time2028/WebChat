package routes

import (
	"github.com/mr-time2028/WebChat/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var RouteRepo *RouteRepository

type RouteRepository struct {
	App *config.App
}

func NewRoutes(a *config.App) {
	RouteRepo = &RouteRepository{
		App: a,
	}
}

func (r *RouteRepository) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Mount("/", r.WsRoutes())
	mux.Mount("/users", r.UserRoutes())

	fileServer := http.FileServer(http.Dir("./web/static/")) // path from root level of the project
	mux.Handle("/web/static/*", http.StripPrefix("/web/static", fileServer))

	return mux
}
