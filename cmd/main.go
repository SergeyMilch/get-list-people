package main

import (
	"log"
	"os"

	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/consumer"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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

	consumer.Start(kafkaBrokers, kafkaTopic, db)
}
