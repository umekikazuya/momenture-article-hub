package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	// Health check endpoint
	http.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check endpoint hit")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Printf("Server starting on port %s...\n", "8080")
	log.Fatal(http.ListenAndServe(":"+"8080", nil))
}
