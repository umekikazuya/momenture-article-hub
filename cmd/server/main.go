package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/umekikazuya/momenture-article-hub/internal/config"
	"github.com/umekikazuya/momenture-article-hub/internal/infrastructure/persistence/postgres"
)

func main() {
	config, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatal("Failed to load configuration")
	}
	// データベース接続。
	db, err := postgres.NewPostgreSQLDB(&config.Database)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	fmt.Println("Database connection established successfully:", db)
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
