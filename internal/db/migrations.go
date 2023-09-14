package db

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func ExecMigrations(db *sqlx.DB) error {
	m, err := migrate.New("file://internal/db/migrations", os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println("Error1")
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("Error2")
		log.Fatal(err)
	}
	return nil
}
