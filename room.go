package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const (
	EventNewMSG = "NEW_MSG_EVENT"
	//EventSendMSG        = "SEND_MSG_EVENT"
	EventNewUser        = "NEW_USER_EVENT"
	EventUserDisconnect = "USER_DISCONNECT_EVENT"
	//EventTyping         = "USER_TYPING_EVENT"
	//EventTypingOut = "USER_TYPING_OUT_EVENT"
	EventAllUsers = "ALL_USERS_EVENT"
)

type Event struct {
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type client struct {
	socket *websocket.Conn
	name   string
}

func newClient(conn *websocket.Conn, name string) *client {
	return &client{
		socket: conn,
		name:   name,
	}
}

func (c *client) forward(event *[]byte) bool {
	if err := c.socket.WriteMessage(websocket.TextMessage, *event); err != nil {
		return false
	}
	return true
}

func (c *client) Read() {
	defer c.socket.Close()
	for {
		_, payload, err := c.socket.ReadMessage()
		if err != nil {
			handleDisconnect(c)
			return
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		switch request.Type {
		case EventNewMSG:
			var data struct {
				To string `json:"to"`
			}
			if err := json.Unmarshal(request.Payload, &data); err != nil {
				log.Printf("error unmarshalling NewMSG payload: %v", err)
				continue
			}

			if user, ok := Users[data.To]; ok {
				// Ensure proper error handling for WriteMessage
				if err := user.socket.WriteMessage(websocket.TextMessage, payload); err != nil {
					log.Printf("error writing message to user: %v", err)
				}
			}
		case EventAllUsers:
			// Handle error for JSON marshalling
			data := struct {
				Type  string              `json:"type"`
				Users map[string]struct{} `json:"users"`
			}{
				Users: getUsersList(),
				Type:  EventAllUsers,
			}
			byts, err := json.Marshal(data)
			if err != nil {
				log.Printf("error marshalling all users data: %v", err)
				continue
			}
			// Ensure proper error handling for WriteMessage
			if err := c.socket.WriteMessage(websocket.TextMessage, byts); err != nil {
				log.Printf("error writing all users data to client: %v", err)
			}
		default:
			// Proper error handling for broadcasting
			if err := broadcastMessage(&payload); err != nil {
				log.Printf("error broadcasting message: %v", err)
			}
		}
	}
}

func broadcastMessage(event *[]byte) error {
	for name, user := range Users {
		if !user.forward(event) {
			delete(Users, name) // Connection Closed
			if data, err := json.Marshal(Event{
				Type: EventUserDisconnect,
				Name: name,
			}); err == nil {
				broadcastMessage(&data)
			}
		}
	}
	return nil
}

func handleDisconnect(c *client) {
	delete(Users, c.name)
	if data, err := json.Marshal(Event{
		Type: EventUserDisconnect,
		Name: c.name,
	}); err == nil {
		// Ensure proper error handling for broadcasting disconnect event
		if err := broadcastMessage(&data); err != nil {
			log.Printf("error broadcasting user disconnect: %v", err)
		}
	}
}

var (
	userMutex sync.RWMutex // RWMutex for synchronization
)

func getUsersList() map[string]struct{} {
	userMutex.RLock()
	defer userMutex.RUnlock()

	userList := make(map[string]struct{})
	for name := range Users {
		userList[name] = struct{}{}
	}
	return userList
}
