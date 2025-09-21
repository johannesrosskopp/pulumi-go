package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "v2 Current time: %s", time.Now().Format(time.RFC3339))
	})

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
