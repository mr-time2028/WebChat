package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mr-time2028/WebChat/internal/helpers"
	"github.com/mr-time2028/WebChat/internal/models"
	"github.com/mr-time2028/WebChat/web/render"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (h *HandlerRepository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index_page.html", &render.TemplateData{})
}

func (h *HandlerRepository) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"name"`
	}

	var responseBody struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	// get user data and json validation
	if validator := helpers.ReadJSON(w, r, &requestBody); !validator.Valid() {
		if err := helpers.ErrorMapJSON(w, validator.Errors); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// create room
	room := &models.Room{
		Name:    requestBody.Name,
		Clients: make(map[string]*models.Client),
	}

	// add room to the database
	roomID, err := h.App.Models.Room.InsertOneRoom(room)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create room, try again", http.StatusInternalServerError)
		return
	}

	// add room to rooms list in hub
	h.App.Hub.Rooms[roomID] = room

	// return response
	responseBody.Error = false
	responseBody.Message = fmt.Sprintf("room with name '%s' successfully created", requestBody.Name)
	if err = helpers.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerRepository) JoinRoom(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token  string `json:"token"`
		RoomID string `json:"room_id"`
	}

	var responseBody struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	// upgrade connection
	conn, err := models.UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = conn.ReadJSON(&requestBody)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		_ = conn.Close()
		return
	}

	claims, err := h.App.Auth.VerifyAuthToken(requestBody.Token)
	if err != nil {
		responseBody.Error = true
		responseBody.Message = "invalid user credentials"
		if err = conn.WriteJSON(responseBody); err != nil {
			log.Println("failed to write json to ws")
		}
		_ = conn.Close()
		return
	}

	userID := claims.Subject
	user, err := h.App.Models.User.GetUserByID(userID)
	if err != nil {
		responseBody.Error = true
		switch err {
		case gorm.ErrRecordNotFound:
			responseBody.Message = "invalid user credentials"
		default:
			responseBody.Message = "internal server error"
		}
		if err = conn.WriteJSON(responseBody); err != nil {
			log.Println("failed to write json to ws")
		}
		_ = conn.Close()
		return
	}

	// register a new client
	client := &models.Client{
		Conn:     conn,
		Hub:      h.App.Hub,
		Message:  make(chan *models.Message),
		ID:       userID,
		RoomID:   requestBody.RoomID,
		Username: user.Username,
	}
	h.App.Hub.Register <- client

	// broadcast join message to all client in that room
	m := &models.Message{
		Sender:   *user,
		SenderID: userID,
		RoomID:   requestBody.RoomID,
		Content:  "A new user has joined the room",
	}
	h.App.Hub.Broadcast <- m

	// start read and write message for client
	go client.ReadMessage()
	go client.WriteMessage()
}

func (h *HandlerRepository) GetRooms(w http.ResponseWriter, r *http.Request) {
	type RoomRes struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var responseBody struct {
		Error bool      `json:"error"`
		Rooms []RoomRes `json:"rooms"`
	}

	rooms := make([]RoomRes, 0)

	for _, r := range h.App.Hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	// return response
	responseBody.Error = false
	responseBody.Rooms = rooms
	if err := helpers.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerRepository) GetRoomClients(w http.ResponseWriter, r *http.Request) {
	type ClientRes struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	var requestBody struct {
		RoomID string `json:"room_id"`
	}

	var responseBody struct {
		Error   bool        `json:"error"`
		Clients []ClientRes `json:"clients"`
	}

	// get user data and json validation
	if validator := helpers.ReadJSON(w, r, &requestBody); !validator.Valid() {
		if err := helpers.ErrorMapJSON(w, validator.Errors); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	roomID := requestBody.RoomID
	clients := make([]ClientRes, 0)

	// TODO: we should return error when no rooms exist with this roomID
	if _, ok := h.App.Hub.Rooms[roomID]; !ok {
		responseBody.Error = true
		responseBody.Clients = clients
		if err := helpers.ErrorStrJSON(w, errors.New("no room exists with this id"), http.StatusBadRequest); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	for _, c := range h.App.Hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	responseBody.Error = false
	responseBody.Clients = clients
	if err := helpers.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
