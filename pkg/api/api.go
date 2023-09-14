package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	ID          uint   `db:"id" json:"id"`
	UserName    string `db:"user_name" json:"user_name"`
	Surname     string `db:"surname" json:"surname"`
	Patronymic  string `db:"patronymic" json:"patronymic"`
	Age         int    `db:"age" json:"age"`
	Gender      string `db:"gender" json:"gender"`
	Nationality string `db:"nationality" json:"nationality"`
}

func GetAge(name string) (uint8, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	age := uint8(data["age"].(float64))
	return age, nil
}

func GetGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	gender := data["gender"].(string)
	return gender, nil
}

// Поскольку api.nationalize.io возвращает массив с некими значениями вероятностей, то в национальность
// запишем наибольшее значение "country_id"
func GetNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	countries, ok := data["country"].([]interface{})
	if !ok || len(countries) == 0 {
		return "", fmt.Errorf("Не удалось получить данные о странах")
	}

	var mostProbableCountryID string
	var maxProbability float64 = 0

	for _, country := range countries {
		countryData, ok := country.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("Не удалось получить данные о стране")
		}

		probability, ok := countryData["probability"].(float64)
		if !ok {
			return "", fmt.Errorf("Не удалось получить вероятность")
		}

		if probability > maxProbability {
			maxProbability = probability
			mostProbableCountryID, ok = countryData["country_id"].(string)
			if !ok {
				return "", fmt.Errorf("Не удалось получить ID страны")
			}
		}
	}

	if mostProbableCountryID == "" {
		return "", fmt.Errorf("Не удалось найти наиболее вероятную страну")
	}

	return mostProbableCountryID, nil
}

func GetPeople(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var people []Person
		err := db.Select(&people, "SELECT * FROM people")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, people)
	}
}

func AddPerson(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var person Person
		err := c.BindJSON(&person)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		_, err = db.Exec("INSERT INTO people (user_name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)",
			person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

func DeletePerson(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM people WHERE id = $1", id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

func UpdatePerson(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var person Person
		err := c.BindJSON(&person)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		_, err = db.Exec("UPDATE people SET user_name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7",
			person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}
