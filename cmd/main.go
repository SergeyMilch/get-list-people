package main

import (
	"log"
	"os"

	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/consumer"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = execMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	consumer.Start(kafkaBrokers, kafkaTopic, db)
}

func execMigrations(db *sqlx.DB) error {
	m, err := migrate.New("file://migrations", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	return nil
}
