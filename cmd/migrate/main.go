package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/awahids/bn-server/configs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	action := flag.String("action", "up", "migration action: up | down | steps | version | force")
	steps := flag.Int("steps", 0, "steps count for action=steps")
	force := flag.Int("force", 0, "target version for action=force")
	path := flag.String("path", "internal/infrastructure/database/migrations", "migration files path")
	flag.Parse()

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	absPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf("failed to resolve migration path: %v", err)
	}

	migrationSource := fmt.Sprintf("file://%s", absPath)
	m, err := migrate.New(migrationSource, cfg.DB.URL())
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	switch *action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "steps":
		if *steps == 0 {
			log.Fatal("--steps cannot be zero for action=steps")
		}
		err = m.Steps(*steps)
	case "version":
		version, dirtyErr := printVersion(m)
		if dirtyErr != nil {
			log.Fatalf("failed to read version: %v", dirtyErr)
		}
		log.Printf("migration version: %d", version)
		return
	case "force":
		if *force <= 0 {
			log.Fatal("--force must be > 0 for action=force")
		}
		err = m.Force(*force)
	default:
		log.Fatalf("unknown action: %s", *action)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration failed: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("no migration changes")
	} else {
		log.Printf("migration action '%s' executed successfully", *action)
	}

	if _, dirtyErr := printVersion(m); dirtyErr != nil {
		log.Fatalf("failed to read migration version: %v", dirtyErr)
	}
}

func printVersion(m *migrate.Migrate) (uint, error) {
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Println("migration version: none")
			return 0, nil
		}
		return 0, err
	}

	if dirty {
		return version, fmt.Errorf("dirty migration state at version %d", version)
	}

	log.Printf("current migration version: %d", version)
	return version, nil
}
