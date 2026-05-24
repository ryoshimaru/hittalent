package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ryoshimaru/hittalent/internal/config"
	"github.com/ryoshimaru/hittalent/internal/database"
)

func main() {
	cfg := config.Load()

	_, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("failed to connect to database : ", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	address := fmt.Sprintf(":%s", cfg.HTTPPort)

	log.Println("server started on", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("server error: ", err)
	}

}
