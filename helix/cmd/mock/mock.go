package main

import (
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/backend"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", backend.MockHandler)

	log.Println("Mock backend running on :8081")

	log.Fatal(http.ListenAndServe(":8081", mux))
}
