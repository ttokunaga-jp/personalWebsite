package firestoredb

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"go.uber.org/fx"
	"google.golang.org/api/option"

	"github.com/takumi/personal-website/internal/config"
)

// NewClient initialises a Firestore client backed by application configuration.
// When the project ID is not configured the caller receives nil, allowing tests
// or local runs without Firestore to fall back to in-memory repositories.
func NewClient(lc fx.Lifecycle, cfg *config.AppConfig) (*firestore.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("firestore client: missing app config")
	}

	fsCfg := cfg.Firestore
	if fsCfg.ProjectID == "" {
		log.Println("firestore client: project id not configured; persistence disabled")
		return nil, nil
	}

	if fsCfg.EmulatorHost != "" {
		// Allow running against the local emulator without needing global env tweaks.
		if err := os.Setenv("FIRESTORE_EMULATOR_HOST", fsCfg.EmulatorHost); err != nil {
			return nil, fmt.Errorf("firestore client: set emulator host: %w", err)
		}
	}

	ctx := context.Background()

	if fsCfg.DatabaseID != "" && fsCfg.DatabaseID != "(default)" {
		log.Printf("firestore client: custom database_id %q requested but NewClientWithConfig is unavailable; falling back to default database", fsCfg.DatabaseID)
	}

	client, err := firestore.NewClient(ctx, fsCfg.ProjectID, clientOptions(fsCfg)...)
	if err != nil {
		return nil, fmt.Errorf("firestore client: create: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

func clientOptions(cfg config.FirestoreConfig) []option.ClientOption {
	// Placeholder for future options (quota project, custom endpoint, etc.).
	// Returning nil would panic in variadic call, so ensure we always return a slice.
	return []option.ClientOption{}
}
