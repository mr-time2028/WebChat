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

	go h.ReadMessage(conn)
}

func (h *HandlerRepository) ReadMessage(conn *websocket.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	var wsRequest models.WsRequest
	for {
		err := conn.ReadJSON(&wsRequest)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		} else {
			client := models.Client{
				Hub:  h.App.Hub,
				Conn: conn,
			}
			wsRequest.Client = &client
			h.App.Hub.RequestChan <- wsRequest
		}
	}
}

func (h *HandlerRepository) WsBroker() {
	for {
		cr := <-h.App.Hub.RequestChan
		switch cr.Action {
		case "auth":
			ok := h.Authenticate(&cr)
			if !ok {
				_ = cr.Client.Conn.Close()
			}
		case "left":
			h.App.Hub.RemoveClient(cr.Client.Conn)
		default:
			wsResponse := &models.WsResponse{
				Error:   true,
				Status:  http.StatusBadRequest,
				Message: "invalid action",
			}
			_ = cr.Client.Conn.WriteJSON(wsResponse)
			_ = cr.Client.Conn.Close()
		}
	}
}

func (h *HandlerRepository) Authenticate(cr *models.WsRequest) bool {
	wsResponse := &models.WsResponse{
		Action: cr.Action,
	}

	// validate auth token
	claims, err := h.App.Auth.VerifyAuthToken(cr.Token)
	if err != nil {
		wsResponse.Error = true
		wsResponse.Status = http.StatusUnauthorized
		wsResponse.Message = "token is invalid or expired"
		_ = cr.Client.Conn.WriteJSON(wsResponse)
		return false
	}

	// get user from database
	parsedUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		wsResponse.Error = true
		wsResponse.Status = http.StatusInternalServerError
		wsResponse.Message = "internal server error"
		_ = cr.Client.Conn.WriteJSON(wsResponse)
		return false
	}

	user, err := h.App.Models.User.GetUserByID(parsedUUID)
	if err != nil {
		wsResponse.Error = true
		wsResponse.Status = http.StatusUnauthorized
		wsResponse.Message = "invalid user credentials"
		_ = cr.Client.Conn.WriteJSON(wsResponse)
		return false
	}

	h.App.Hub.AddClient(cr.Client.Conn, user)
	wsResponse.ConnectedUsers = h.App.Hub.GetConnectedUsers()
	wsResponse.Error = false
	wsResponse.Status = http.StatusOK
	wsResponse.Message = "user authenticated successfully"
	_ = cr.Client.Conn.WriteJSON(wsResponse)
	return true
}
