package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/undg/go-prapi/buildinfo"
	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
	"github.com/undg/go-prapi/ws"
)

// @TODO (undg) 2024-10-06: different port for dev and production

func startServer(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api":
			json.RenderSchemaJSON(w, r)
		case "/api/v1/ws":
			ws.HandleWebSocket(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	fs := http.FileServer(http.Dir("/tmp/bin/pr-web/dist"))
	mux.Handle("/", fs)
}

func main() {
	ip := utils.GetLocalIP()
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

	fmt.Println("\n🔥 Igniting server on ws://" + ip + utils.PORT + "\n")

	mux := http.NewServeMux()

	startServer(mux)

	go pactl.ListenForChanges(ws.BroadcastUpdates)

	err := http.ListenAndServe(utils.PORT, mux)
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}
