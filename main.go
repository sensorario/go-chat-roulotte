package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	// create new hub
	flag.Parse()
	hub := newHub()
	go hub.run()

	// handle endpoints
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/config.js", serveConfig)
	http.HandleFunc("/game.js", serveJs)
	http.HandleFunc("/game.css", serveCss)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// log if error
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
