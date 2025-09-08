package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("gateway init: %v", err)
	}

	a.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ gateway shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Stop(ctx); err != nil {
		log.Printf("gateway stop error: %v", err)
	}
	log.Println("✅ gateway bye")
}
