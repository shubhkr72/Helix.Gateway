package backend

import (
	"encoding/json"
	"net/http"
	"os"
)

type Response struct {
	Instance string            `json:"instance"`
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Headers  map[string]string `json:"headers"`
}

func MockHandler(w http.ResponseWriter, r *http.Request) {

	instance := os.Getenv("INSTANCE_NAME")

	if instance == "" {
		instance = "mock-1"
	}

	resp := Response{
		Instance: instance,
		Method:   r.Method,
		Path:     r.URL.Path,
		Headers: map[string]string{
			"User-Agent":    r.Header.Get("User-Agent"),
			"Authorization": r.Header.Get("Authorization"),
			"Content-Type":  r.Header.Get("Content-Type"),
		},
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}
