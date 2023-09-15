package api

import (
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
		err := c.ShouldBindJSON(&person)
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
		err := c.ShouldBindJSON(&person)
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
