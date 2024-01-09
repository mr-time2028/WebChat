package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Client contains websocket connection instance for each user
type Client struct {
	*websocket.Conn
}

var Clients = make(map[Client]string)

var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WsRequest contains what the clients send to the server
type WsRequest struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	ConnectedUsers []string `json:"connected_users"`
}

// WsResponse contains what the server sends to the clients
type WsResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}
