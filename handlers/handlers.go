package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/mr-time2028/WebChat/helpers"
	"github.com/mr-time2028/WebChat/render"
	"github.com/mr-time2028/WebChat/server"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index_page.html", &render.TemplateData{})
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := server.UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		helpers.ErrorStrJSON(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	log.Println("client connected!")

	ws.WriteJSON(&server.WsResponse{
		Action: "success",
		Message: "this is worked!",
	})
}
