package main

import (
	"log"
	"net/http"
)

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
