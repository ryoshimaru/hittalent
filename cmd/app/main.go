package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ryoshimaru/hittalent/internal/config"
	"github.com/ryoshimaru/hittalent/internal/database"
	"github.com/ryoshimaru/hittalent/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	handler := router.New(db)

	address := fmt.Sprintf(":%s", cfg.HTTPPort)

	log.Println("server started on", address)

	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatal("server error: ", err)
	}
}
