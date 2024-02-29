package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mr-time2028/WebChat/internal/models"
	"github.com/mr-time2028/WebChat/web/render"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *HandlerRepository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index_page.html", &render.TemplateData{})
}

func (h *HandlerRepository) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade connection
	conn, err := models.UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal config error", http.StatusInternalServerError)
		return
	}

	go h.ParseWsRequest(conn)
}

func (h *HandlerRepository) ParseWsRequest(conn *websocket.Conn) {
	// TODO: should log it to mongo or write to a file
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", r)
		}
		conn.Close()
		// TODO: i should remove this conn.Close() because i wanna keep connection alive
	}()

	var wsRequest models.WsRequest
	for {
		if err := conn.ReadJSON(&wsRequest); err != nil {
			log.Println(err)
			return
		} else {
			client := models.Client{Conn: conn}
			wsRequest.Client = &client
			h.App.Hub.RequestChan <- wsRequest
		}
	}
}

func (h *HandlerRepository) WsBroker() {
	var response models.WsResponse

	for {
		clientRequest := <-h.App.Hub.RequestChan
		switch clientRequest.Action {
		case "auth":
			user, err := h.Authenticate(&clientRequest)
			if err != nil {
				switch err {
				case models.ErrNoAuthHeader:
					// invalid auth header error
					log.Println(err)
				case gorm.ErrRecordNotFound:
					// no user found error
					log.Println(err)
				default:
					// config error
					log.Println(err)
				}
				return
			}
			h.App.Hub.Clients[clientRequest.Client] = user.Username
			response.Message = "user logged in"
			response.Action = "auth"
			err = clientRequest.Client.WriteJSON(response)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (h *HandlerRepository) Authenticate(clientRequest *models.WsRequest) (*models.User, error) {
	// validate auth token
	claims, err := h.App.Auth.VerifyAuthToken(clientRequest.Authorization)
	if err != nil {
		return nil, err
	}

	// get user from database
	parsedUUID, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, err
	}

	user, err := h.App.Models.User.GetUserByID(parsedUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// TODO: make an error engine helper functions for ws connections
