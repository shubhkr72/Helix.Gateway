package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Service string `json:"service"`
	Message string `json:"message"`
	Method  string `json:"method"`
	Path    string `json:"path"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "users",
	})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Service: "users",
		Message: "User service is running",
		Method:  r.Method,
		Path:    r.URL.Path,
	}

	writeJSON(w, http.StatusOK, resp)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "users",
		"user": map[string]any{
			"id":    "123",
			"name":  "Shubham",
			"email": "user@example.com",
		},
	})
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler)

	// Because gateway uses strip_prefix: true,
	// GET /users -> backend receives "/"
	mux.HandleFunc("/", usersHandler)

	mux.HandleFunc("/profile", profileHandler)

	server := &http.Server{
		Addr:    ":9002",
		Handler: mux,
	}

	log.Println("Users Service listening on :9002")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}