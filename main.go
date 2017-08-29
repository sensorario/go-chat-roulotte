// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

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
	http.ServeFile(w, r, "home.html")
}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "game.js")
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "game.css")
}

func serveConfig(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "config.js")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/config.js", serveConfig)
	http.HandleFunc("/game.js", serveJs)
	http.HandleFunc("/game.css", serveCss)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(hub, w, r) })

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
