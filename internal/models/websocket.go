package models

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var (
	UpgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// Client contains websocket connection instance for each user
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
}

// Hub maintains all things about websocket connection and clients
type Hub struct {
	Clients      map[*websocket.Conn]string
	RequestChan  chan WsRequest
	ResponseChan chan WsResponse
}

func (h *Hub) AddClient(clientConn *websocket.Conn, user *User) {
	// we should prevent add duplicate username (because of different device connecting)
	isExists := false
	for _, username := range h.Clients {
		if username == user.Username {
			isExists = true
			break
		}
	}

	if !isExists {
		h.Clients[clientConn] = user.Username
	} else {
		var wsResponse WsResponse
		wsResponse.Error = true
		wsResponse.Message = "client already connected with another device"
		_ = clientConn.WriteJSON(wsResponse)
		_ = clientConn.Close()
	}
}

func (h *Hub) RemoveClient(clientConn *websocket.Conn) {
	if _, ok := h.Clients[clientConn]; ok {
		_ = clientConn.Close()
		delete(h.Clients, clientConn)
	}
}

func (h *Hub) GetConnectedUsers() []string {
	var userList []string
	for _, x := range h.Clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)
	return userList
}

func (h *Hub) BroadcastToAll(response WsResponse) {
	for clientConn := range h.Clients {
		err := clientConn.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			h.RemoveClient(clientConn)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*websocket.Conn]string),
		RequestChan:  make(chan WsRequest),
		ResponseChan: make(chan WsResponse),
	}
}

// WsRequest contains what the clients send to the config
type WsRequest struct {
	Token       string  `json:"authorization"`
	Action      string  `json:"action"`
	Message     string  `json:"message"`
	MessageType string  `json:"message_type"`
	Client      *Client `json:"-"`
}

// WsResponse contains what the config sends to the clients
type WsResponse struct {
	Error          bool     `json:"error"`
	Status         int      `json:"status"`
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}
