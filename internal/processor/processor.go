package processor

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/SergeyMilch/get-list-people-effective-mobile/pkg/api"
	"github.com/jmoiron/sqlx"
)

type PersonInfo struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         uint8
	Gender      string
	Nationality string
}

func ProcessFIO(msg *sarama.ConsumerMessage, db *sqlx.DB) error {
	var data map[string]string
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		return fmt.Errorf("Ошибка разбора JSON: %s", err)
	}

	name, ok := data["name"]
	if !ok {
		return fmt.Errorf("Отсутствует поле 'name'")
	}

	surname, ok := data["surname"]
	if !ok {
		return fmt.Errorf("Отсутствует поле 'surname'")
	}

	patronymic, ok := data["patronymic"]
	if !ok {
		patronymic = "" // Необязательное поле
	}

	age, err := api.GetAge(name)
	if err != nil {
		log.Printf("Ошибка получения возраста: %s\n", err)
	}

	gender, err := api.GetGender(name)
	if err != nil {
		log.Printf("Ошибка получения пола: %s\n", err)
	}

	nationality, err := api.GetNationality(name)
	if err != nil {
		log.Printf("Ошибка получения национальности: %s\n", err)
	}

	personInfo := PersonInfo{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	// TODO здесь запись в базу данных
	_, err = db.NamedExec(`INSERT INTO people (user_name, surname, patronymic, age, gender, nationality)
		VALUES (:name, :surname, :patronymic, :age, :gender, :nationality)`, personInfo)
	if err != nil {
		return fmt.Errorf("Ошибка записи в базу данных: %s", err)
	}

	fmt.Printf("Обработано: %+v\n", personInfo)

	return nil
}
