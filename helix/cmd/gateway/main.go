package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shubhkr72/helix/internal/handlers"
	"github.com/shubhkr72/helix/internal/proxy"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("%s %s", r.Method, r.URL.String())

		next.ServeHTTP(w, r)
	})
}

func main() {

	target := os.Getenv("BACKEND")
	if target == "" {
		target = "http://localhost:9000"
	}

	p, err := proxy.New(target)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlers.Health)
	mux.HandleFunc("/readyz", handlers.Ready)

	// local admin routes
	mux.Handle("/admin/", http.NotFoundHandler())

	// everything else goes to backend
	mux.Handle("/", p)

	server := &http.Server{

		Addr:              ":8080",
		Handler:           logging(mux),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Println("Gateway listening on :8080")

	log.Fatal(server.ListenAndServe())
}