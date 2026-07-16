package main

import (
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/config"
	"github.com/shubhkr72/helix/internal/handlers"
	"github.com/shubhkr72/helix/internal/gateway"
)

func main() {

	cfg, err := config.Load("configs/gateway.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Listening on :%d", cfg.Server.Port)
	gateway.PrintBanner(cfg)

	http.ListenAndServe(
		":"+"8080",
		&handlers.Gateway{
			Config: cfg,
		},
	)
}
