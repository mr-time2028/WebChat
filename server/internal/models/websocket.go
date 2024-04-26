package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	Conn     *websocket.Conn
	Hub      *Hub
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

func (c *Client) ReadMessage() {
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content: string(m),
			RoomID:  c.RoomID,
		}
		c.Hub.Broadcast <- msg
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		_ = c.Conn.WriteJSON(message)
	}
}

// Hub maintains all things about websocket connection and clients
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// check room exists
			if _, ok := h.Rooms[client.RoomID]; ok {
				// add client to the room
				r := h.Rooms[client.RoomID] // specify which room the client request to add to it

				if r.Clients == nil {
					// Initialize the Clients map because a client want to join to this room
					r.Clients = make(map[string]*Client)
				}

				if _, ok = r.Clients[client.ID]; !ok {
					r.Clients[client.ID] = client // add client to the room if client not in the room
				}
			}
		case client := <-h.Unregister:
			// check if room exists
			if _, ok := h.Rooms[client.RoomID]; ok {
				// delete client from room (should search in all rooms in hub and find client)
				if _, ok = h.Rooms[client.RoomID].Clients[client.ID]; ok {
					delete(h.Rooms[client.RoomID].Clients, client.ID)
					close(client.Message)

					// broadcast a message that the client has left the room
					if len(h.Rooms[client.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  fmt.Sprintf("user %s left the chat", client.Username),
							SenderID: client.Username,
							RoomID:   client.RoomID,
						}
					}

					// Set Clients map to nil if it's empty
					if len(h.Rooms[client.RoomID].Clients) == 0 {
						h.Rooms[client.RoomID].Clients = nil
					}
				}
			}
		case message := <-h.Broadcast:
			// check if room exists
			if _, ok := h.Rooms[message.RoomID]; ok {
				// broadcast message to all clients in the room
				for _, client := range h.Rooms[message.RoomID].Clients {
					client.Message <- message
				}
			}
		}
	}
}

//func (h *Hub) AddClient(clientConn *websocket.Conn, user *User) {
//	// we should prevent add duplicate username (because of different device connecting)
//	isExists := false
//	for _, username := range h.Clients {
//		if username == user.Username {
//			isExists = true
//			break
//		}
//	}
//
//	if !isExists {
//		h.Clients[clientConn] = user.Username
//	} else {
//		var wsResponse WsResponse
//		wsResponse.Error = true
//		wsResponse.Message = "client already connected with another device"
//		_ = clientConn.WriteJSON(wsResponse)
//		_ = clientConn.Close()
//	}
//}
//
//func (h *Hub) RemoveClient(clientConn *websocket.Conn) {
//	if _, ok := h.Clients[clientConn]; ok {
//		_ = clientConn.Close()
//		delete(h.Clients, clientConn)
//	}
//}
//
//func (h *Hub) GetConnectedUsers() []string {
//	var userList []string
//	for _, x := range h.Clients {
//		if x != "" {
//			userList = append(userList, x)
//		}
//	}
//	sort.Strings(userList)
//	return userList
//}
//
//func (h *Hub) BroadcastToAll(response WsResponse) {
//	for clientConn := range h.Clients {
//		err := clientConn.WriteJSON(response)
//		if err != nil {
//			log.Println("websocket err")
//			h.RemoveClient(clientConn)
//		}
//	}
//}

type WsAuthPayload struct {
	Token string `json:"token"`
}

// WsRequest contains what the clients send to the config
type WsRequest struct {
	//Token       string  `json:"token"`
	Action        string        `json:"action"`
	WsAuthPayload WsAuthPayload `json:"auth"`
	//Message     string  `json:"message"`
	//MessageType string  `json:"message_type"`
	Client *Client `json:"-"`
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
