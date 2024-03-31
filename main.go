package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8448"

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home Page hit by clien")

	fmt.Fprintf(w, "Welcome on 'home' endpoint")
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
}

func main() {
	fmt.Println("Starting server")

	setupRoutes()

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
