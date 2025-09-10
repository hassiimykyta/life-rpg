package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hassiimykyta/life-rpg/apps/notification-svc/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("app.New: %v", err)
	}

	if err := a.Start(); err != nil {
		log.Fatalf("app.Start: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("shutdown signal: %v", sig)
	case err := <-a.ErrChan():
		if err != nil {
			log.Printf("background error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = a.Stop(ctx)
}
