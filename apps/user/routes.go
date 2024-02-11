package user

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes() http.Handler {
	mux := chi.NewRouter()
	return mux
}
