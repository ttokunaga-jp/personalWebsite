package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/takumi/personal-website/internal/app"
)

func main() {
	root := app.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := root.Start(ctx); err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	<-ctx.Done()

	if err := root.Stop(context.Background()); err != nil {
		log.Printf("graceful shutdown encountered error: %v", err)
	}

	os.Exit(0)
}
