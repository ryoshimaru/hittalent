package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql database: ", err)
	}

	handler := router.New(db)

	address := fmt.Sprintf(":%s", cfg.HTTPPort)

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	go func() {
		log.Println("server started on", address)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatal("failed to close database connection: ", err)
	}

	log.Println("server stopped")
}
