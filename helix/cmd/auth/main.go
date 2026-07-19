package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shubhkr72/helix/internal/authconfig"
	"github.com/shubhkr72/helix/internal/database"
)

func main() {

	fmt.Println("Loading configuration...")

	cfg, err := authconfig.Load("configs/auth.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to PostgreSQL...")

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected successfully.")

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Server.Port),
	}

	fmt.Printf("Auth Service listening on :%d\n", cfg.Server.Port)

	log.Fatal(server.ListenAndServe())
}