package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/takumi/personal-website/internal/app"
)

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		// .env が無い場合は何もしない
		return
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("failed to close env file: %v", cerr)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 囲み文字を削除
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			// Respect already-specified environment variables so CLI overrides are honoured.
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			log.Printf("failed to set env var %s: %v", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading env file: %v", err)
	}
}

func main() {
	loadEnvFile(".env")
	loadEnvFile("../.env")

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
