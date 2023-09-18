package consumer

import (
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/db"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/processor"
	"github.com/jmoiron/sqlx"
)

func Start(brokers, topic string, db db.Database) {
	dbConn, ok := db.(*sqlx.DB)
	if !ok {
		log.Fatal("Ошибка преобразования db.Database в *sqlx.DB")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{brokers}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Ошибка закрытия потребителя: %s\n", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Ошибка закрытия раздела потребителя: %s\n", err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			err := processor.ProcessFIO(msg, dbConn)
			if err != nil {
				log.Printf("Ошибка обработки сообщения: %s\n", err)
				sendToFailedTopic(brokers, msg.Value)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Ошибка: %s\n", err.Error())
		}
	}
}

func sendToFailedTopic(brokers string, msg []byte) {
	producer, err := sarama.NewAsyncProducer([]string{brokers}, nil)
	if err != nil {
		log.Printf("Ошибка создания продюсера: %s\n", err)
		return
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Ошибка закрытия продюсера: %s\n", err)
		}
	}()

	failedMsg := &sarama.ProducerMessage{
		Topic: os.Getenv("KAFKA_FAILED"),
		Value: sarama.StringEncoder(msg),
	}

	producer.Input() <- failedMsg
}
