package main

import (
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/handlers"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlers.Health)
	mux.HandleFunc("/readyz", handlers.Ready)

	log.Println("Gateway running on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}