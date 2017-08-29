package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/config.js", serveConfig)
	http.HandleFunc("/game.js", serveJs)
	http.HandleFunc("/game.css", serveCss)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "client/home.html")
}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/game.js")
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/game.css")
}

func serveConfig(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/config.js")
}
