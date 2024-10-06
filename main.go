package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/buildinfo"
)

// @TODO (undg) 2024-10-06: different port for dev and production

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
	b := buildinfo.Get()
	fmt.Println("\t* GitVersion:\t", b.GitVersion)
	fmt.Println("\t* GitCommit:\t", b.GitCommit)
	fmt.Println("\t* BuildDate:\t", b.BuildDate)
	fmt.Println("\t* Compiler:\t", b.Compiler)
	fmt.Println("\t* Platform:\t", b.Platform)
	fmt.Println("\t* GoVersion:\t", b.GoVersion)

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
