package processor

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func ProcessFIO(msg *sarama.ConsumerMessage) {
	var data map[string]string
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		log.Printf("Ошибка разбора JSON: %s\n", err)
		return
	}

	name, ok := data["name"]
	if !ok {
		log.Println("Отсутствует поле 'name'")
		return
	}

	surname, ok := data["surname"]
	if !ok {
		log.Println("Отсутствует поле 'surname'")
		return
	}

	patronymic, ok := data["patronymic"]
	if !ok {
		patronymic = "" // Необязательное поле
	}

	// TODO здесь запись в базу данных

	fmt.Printf("Обработано: %s %s %s\n", name, surname, patronymic)
}
