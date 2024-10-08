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

	fs := http.FileServer(http.Dir("/tmp/bin/pr-web/dist"))
	mux.Handle("/", fs)
}

func main() {
	ip := getLocalIP()
	b := buildinfo.Get()

	fmt.Print(`
┌───────────────────────────────────────────────────┐
│                     GO-PRAPI                      │
├───────────────────────────────────────────────────┤
  GitVersion: `, b.GitVersion, `
  GitCommit:  `, b.GitCommit, `
  BuildDate:  `, b.BuildDate, `
  Compiler:   `, b.Compiler, `
  Platform:   `, b.Platform, `
  GoVersion:  `, b.GoVersion, `
└───────────────────────────────────────────────────┘
`)

	fmt.Println("\n🔥 Igniting server on ws://" + ip + PORT + "\n")

	mux := http.NewServeMux()

	startServer(mux)

	go broadcastUpdates()

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}
