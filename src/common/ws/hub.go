package ws

import (
	"TTCS/src/common/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Hub struct {
	Clients    map[string]map[uuid.UUID]*Client // Lưu trữ theo route và ID client
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]map[uuid.UUID]*Client),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.RegisterClient(client)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.Route][client.ID]; ok {
				delete(h.Clients[client.Route], client.ID)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			h.BroadcastMessageToRoute(message.Type, message)
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	log.Info(context.Background(), "client with id: %s is registed route %s", client.ID, client.Route)
	if _, exists := h.Clients[client.Route]; !exists {
		h.Clients[client.Route] = make(map[uuid.UUID]*Client)
	}
	h.Clients[client.Route][client.ID] = client
}

func (h *Hub) BroadcastMessageToRoute(route string, message Message) {
	if clients, exists := h.Clients[route]; exists {
		for _, client := range clients {
			data, _ := json.Marshal(message)
			client.Send <- data
		}
	}
}

func (h *Hub) SendMessageToClient(clientID uuid.UUID, message Message) error {
	if clients, exists := h.Clients[message.Type]; exists {
		client, ok := clients[clientID]
		if !ok {
			return fmt.Errorf("client with ID %s not found in route %s", clientID, message.Type)
		}

		data, err := json.Marshal(message)
		if err != nil {
			return err
		}
		client.Send <- data
		return nil
	}
	return fmt.Errorf("route %s not found", message.Type)
}
