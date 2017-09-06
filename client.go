package main

import (
	"bytes"
	//"fmt"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func (c *Client) readPump() {

	// defer unregistration
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// init connection
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// infinite loop
	for {

		// Read connection messages
		log.Println("Aspetto che qualcuno invii un messaggio")
		_, message, err := c.conn.ReadMessage()
		log.Println("Leggo il messaggio")

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Println("Trasformo il messaggio")
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// send everything!!!
		log.Println("Invio il messaggio a tutti")
		c.hub.broadcast <- message

	}

}

// Write ...
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {

		// Whenever a message is sent
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}

// In pratica quando ci si connette bisogna serviere il web service
// che non fa altro che registrare questo client nell'hub
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println(" ::: web service ::: ")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// New client created
	client := &Client{
		hub:      hub,
		opponent: nil,
		conn:     conn,
		send:     make(chan []byte, 256),
	}

	// Il client viene registrato nell'hub
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
