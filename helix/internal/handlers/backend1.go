package handlers

import (
	"encoding/json"
	"net/http"
)

func Backend1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "Backend 1 is up & running",
	})
}
