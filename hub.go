package main

import (
	"log"
)

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {

		// register new client
		case client := <-h.register:
			log.Println("register")
			h.clients[client] = true
			// todo assign a room to this client
			// hub.client[questo client] = true

		// unregister new client
		// this method is called with defer
		case client := <-h.unregister:
			log.Println("unregister")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		// broadcast message
		// Quando si riceve un messaggio broadcast
		case message := <-h.broadcast:
			log.Println("broadcast")
			// all clients connected
			for client := range h.clients {
				select {
				// Lo si rispedisce a tutti i client
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}
