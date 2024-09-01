package main

import (
	"fmt"
	"net/http"
)

func apiDoc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome on 'api' endpoint documentation\n")
	fmt.Fprintf(w, "🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻🧻")
}
