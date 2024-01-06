package server

import (
	"github.com/gorilla/websocket"
)

// Client contains websocket connection instance for each user
type Client struct {
	conn *websocket.Conn
}

// clients = make(chan[*WebSocketConnection]string)

var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
