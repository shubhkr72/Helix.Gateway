package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shubhkr72/helix/internal/auth"
	"github.com/shubhkr72/helix/internal/authconfig"
	"github.com/shubhkr72/helix/internal/database"
	"github.com/shubhkr72/helix/internal/jwt"
)

func main() {

	cfg, err := authconfig.Load("configs/auth.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := auth.NewPostgresRepository(db)

	jwtMgr, err := jwt.NewManager(
		"keys/private.pem",
		"keys/public.pem",
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		time.Duration(cfg.JWT.Expiry)*time.Minute,
	)
	if err != nil {
		log.Fatal(err)
	}

	service := auth.NewService(repo, jwtMgr)

	handler := auth.NewHandler(service)

	mux := http.NewServeMux()

	auth.RegisterRoutes(mux, handler)

	log.Println("Auth Service listening on :9001")

	log.Fatal(http.ListenAndServe(":9001", mux))
}