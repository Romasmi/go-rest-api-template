package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (db *DbConnection) RunMigrations(direction string) error {

	m, err := migrate.New("file://migrations", db.Config.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			fmt.Printf("error while closing DB connection: %v", err)
		}
	}(m)

	if direction == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to run migrations up: %w", err)
		}
		log.Println("Migrations up completed successfully")
	} else if direction == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to run migrations down: %w", err)
		}
		log.Println("Migrations down completed successfully")
	} else {
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	return nil
}

func MigrateToVersion(version uint) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/go_rest_api?sslmode=disable"
	}

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Migrate(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	log.Printf("Migration to version %d completed successfully", version)
	return nil
}
