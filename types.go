package main

import (
	"github.com/gorilla/websocket"
)

// Un Hub contiene una mappa di Client
// un Hub in cui sia possibile registrarsi, deregistrarsi, ...
// mandare messaggi broadcast e registrare una marea di client
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// un client punta ad un Hub
// un client al quale si possono inviare messaggi
// un cliente ha anche un opponent
type Client struct {
	hub      *Hub
	opponent *Client
	conn     *websocket.Conn
	send     chan []byte
}
