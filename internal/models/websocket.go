package models

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// Client contains websocket connection instance for each user
type Client struct {
	*websocket.Conn
}

var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Hub maintains all things about websocket connection and clients
type Hub struct {
	Clients      map[*Client]string
	RequestChan  chan WsRequest
	ResponseChan chan WsResponse
}

func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]string),
		RequestChan:  make(chan WsRequest),
		ResponseChan: make(chan WsResponse),
	}
}

// WsRequest contains what the clients send to the config
type WsRequest struct {
	Authorization string  `json:"-"`
	Action        string  `json:"action"`
	Username      string  `json:"username"`
	Message       string  `json:"message"`
	MessageType   string  `json:"message_type"`
	Client        *Client `json:"-"`
}

// WsResponse contains what the config sends to the clients
type WsResponse struct {
	Error          bool     `json:"error"`
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}
