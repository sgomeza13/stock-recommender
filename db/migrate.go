package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sgomeza13/stock-recommender/utils"
)

// RunMigrations applies database migrations
func RunMigrations() {
	migrationsPath := "file://db/migrations"

	m, err := migrate.New(migrationsPath, utils.GetDSN(true))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migrations: %v", err)
	}

	fmt.Println(" Migrations applied successfully!")
}
