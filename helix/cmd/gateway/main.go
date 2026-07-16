package main

import (
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/handlers"
)

func main() {

	cfg, err := config.Load("configs/gateway.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on :%d", cfg.Server.Port)

	http.ListenAndServe(
		":"+"8080",
		&handlers.Gateway{
			Config: cfg,
		},
	)
}
