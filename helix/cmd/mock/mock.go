package main

import (
	"encoding/json"
	"fmt"
	"github.com/shubhkr72/helix/internal/handlers"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Instance string              `json:"instance"`
	Method   string              `json:"method"`
	Path     string              `json:"path"`
	Query    string              `json:"query"`
	Headers  map[string][]string `json:"headers"`
}

func main() {

	name := os.Getenv("INSTANCE")
	if name == "" {
		name = "backend-1"
	}
	http.HandleFunc("/backend", handlers.Backend1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.Method, r.URL)

		resp := Response{
			Instance: name,
			Method:   r.Method,
			Path:     r.URL.Path,
			Query:    r.URL.RawQuery,
			Headers: map[string][]string{
				"X-Forwarded-For":   r.Header["X-Forwarded-For"],
				"X-Forwarded-Proto": r.Header["X-Forwarded-Proto"],
				"User-Agent":        r.Header["User-Agent"],
			},
		}

		json.NewEncoder(w).Encode(resp)
	})

	log.Println("mock backend :9000")

	log.Fatal(http.ListenAndServe(":9000", nil))
}
