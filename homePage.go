package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home Page hit by clien")

	fmt.Fprintf(w, "Welcome on 'home' endpoint")
}
