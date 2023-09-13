package consumer

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/processor"
)

func Start(brokers, topic string) {
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
			processor.ProcessFIO(msg)
		case err := <-partitionConsumer.Errors():
			log.Printf("Ошибка: %s\n", err.Error())
		}
	}
}
