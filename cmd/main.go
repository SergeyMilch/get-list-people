package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/consumer"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/db"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/router"
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

	dbConn, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println("Error5")
		log.Fatal(err)
	}
	defer dbConn.Close()

	err = db.ExecMigrations(dbConn)
	if err != nil {
		fmt.Println("Error3")
		log.Fatal(err)
	}

	r := router.NewRouter(dbConn)

	err = r.Run(":8000")
	if err != nil {
		fmt.Println("Error4")
		log.Fatal(err)
	}

	consumer.Start(kafkaBrokers, kafkaTopic, dbConn)
}
