package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const PORT = ":8448"
const DEBUG = false

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{}

func startServer(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api":
			RenderSchemaJSON(w, r)
		case "/api/v1/ws":
			HandleWebSocket(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)
}

func main() {
	ip := getLocalIP()

	fmt.Println("Listening on ws://" + ip + PORT)
	mux := http.NewServeMux()

	startServer(mux)

	go broadcastUpdates()

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}
