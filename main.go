package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8448"

func setupRoutes(mux *http.ServeMux) {
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
	fmt.Println("Listening on http://localhost" + PORT)
	mux := http.NewServeMux()
	setupRoutes(mux)

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}
