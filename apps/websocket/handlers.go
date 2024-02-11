package websocket

import (
	"errors"
	"github.com/mr-time2028/WebChat/helpers"
	"github.com/mr-time2028/WebChat/models"
	"github.com/mr-time2028/WebChat/web/render"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index_page.html", &render.TemplateData{})
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := models.UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		helpers.ErrorStrJSON(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	log.Println("client connected!")

	// default username for each new connected user is ""
	// TODO: we should handle it with session (decide front side or back side)
	cfg.Clients[models.Client{Conn: ws}] = ""

	go listenForWs()
}

func listenForWs() {
	// TODO: should log it to mongo or write to a file
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", r)
		}
	}()

	// TODO: we should handle user request by a function and switch case, but we should before this get username
}
