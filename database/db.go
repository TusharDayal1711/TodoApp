package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sqlx.DB

func Init(connectionStr string) {
	var err error
	DB, err = sqlx.Connect("postgres", connectionStr)
	//+v for Detailed value with stack trace
	if err != nil {
		log.Fatalf("Cannot connect to Postgres: %+v", err)
	}
	fmt.Println("Connected to PostgreSQL...")

	if err := migrateUp(DB); err != nil {
		log.Fatalf("Migration failed: %+v", err)
	}
}

func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://database/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return err
	}

	fmt.Println("Migration complete.")
	return nil
}
