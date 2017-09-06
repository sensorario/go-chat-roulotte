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

type Message struct {
	typology string
}

func (h *Hub) run() {
	for {
		select {

		// register new client
		case client := <-h.register:
			log.Println("register")
			h.clients[client] = true
			for c := range h.clients {
				if c != client {
					if c.opponent == nil {
						c.opponent = client
						client.opponent = c
						dat := []byte(`{"foo":"bar"}`)
						c.send <- dat
						client.send <- dat
						log.Println("Dovrei aver avvisato i client")
					}
				}
			}

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
