package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8448"

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Starting server")

	setupRoutes()

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
