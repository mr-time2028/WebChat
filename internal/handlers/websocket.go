package handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mr-time2028/WebChat/internal/models"
	"github.com/mr-time2028/WebChat/web/render"
	"log"
	"net/http"
)

func (h *HandlerRepository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index_page.html", &render.TemplateData{})
}

func (h *HandlerRepository) ServeWs(w http.ResponseWriter, r *http.Request) {
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
	defer func() {
		_ = conn.Close()
	}()

	var wsRequest models.WsRequest
	for {
		err := conn.ReadJSON(&wsRequest)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message: ", err)
			}
			break
		} else {
			client := models.Client{Conn: conn}
			wsRequest.Client = &client
			h.App.Hub.RequestChan <- wsRequest
		}
	}
}

func (h *HandlerRepository) WsBroker() {
	for {
		clientRequest := <-h.App.Hub.RequestChan
		switch clientRequest.Action {
		case "auth":
			h.Authenticate(&clientRequest)
		default:
			wsResponse := &models.WsResponse{
				Error:   true,
				Status:  http.StatusBadRequest,
				Message: "invalid action",
			}
			_ = clientRequest.Client.WriteJSON(wsResponse)
			_ = clientRequest.Client.Close()
		}
	}
}

func (h *HandlerRepository) Authenticate(clientRequest *models.WsRequest) {
	wsResponse := &models.WsResponse{
		Action: clientRequest.Action,
	}

	// validate auth token
	claims, err := h.App.Auth.VerifyAuthToken(clientRequest.Authorization)
	if err != nil {
		wsResponse.Error = true
		wsResponse.Status = http.StatusUnauthorized
		wsResponse.Message = "token is invalid or expired"
		_ = clientRequest.Client.WriteJSON(wsResponse)
		_ = clientRequest.Client.Close()
		return
	}

	// get user from database
	parsedUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		wsResponse.Error = true
		wsResponse.Status = http.StatusInternalServerError
		wsResponse.Message = "internal server error"
		_ = clientRequest.Client.WriteJSON(wsResponse)
		_ = clientRequest.Client.Close()
		return
	}

	user, err := h.App.Models.User.GetUserByID(parsedUUID)
	if err != nil {
		wsResponse.Error = true
		wsResponse.Status = http.StatusUnauthorized
		wsResponse.Message = "invalid user credentials"
		_ = clientRequest.Client.WriteJSON(wsResponse)
		_ = clientRequest.Client.Close()
		return
	}

	h.App.Hub.AddClient(clientRequest.Client, user)
	wsResponse.Error = false
	wsResponse.Status = http.StatusOK
	wsResponse.Message = "user authenticated successfully"
	_ = clientRequest.Client.WriteJSON(wsResponse)
}
