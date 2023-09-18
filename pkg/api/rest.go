package api

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/db"
	"github.com/SergeyMilch/get-list-people-effective-mobile/internal/db/redisdb"
	"github.com/SergeyMilch/get-list-people-effective-mobile/pkg/logger"
	"github.com/gin-gonic/gin"
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

func GetPeople(db db.Database, rdb redisdb.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		cacheKey := "all_people"

		result, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var people []Person
			json.Unmarshal([]byte(result), &people)
			c.JSON(200, people)
			return
		}

		var people []Person
		err = db.Select(&people, "SELECT * FROM people")
		if err != nil {
			logger.Error("Ошибка при получении пользователей из базы данных: ", err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		jsonData, _ := json.Marshal(people)
		err = rdb.Set(ctx, cacheKey, jsonData, time.Minute*10).Err()

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, people)
	}
}

func GetPersonByID(db db.Database, rdb redisdb.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		cacheKey := "person_" + id // Можно сделать хеш, например: keyCache := MD5(cacheKey)

		ctx := context.Background()
		result, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var person Person
			json.Unmarshal([]byte(result), &person)
			c.JSON(200, person)
			return
		}

		var person Person
		err = db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
		if err != nil {
			logger.Warn("Пользователь не найден: ", err.Error())
			c.JSON(404, gin.H{"error": "Пользователь не найден"})
			return
		}

		jsonData, _ := json.Marshal(person)
		err = rdb.Set(ctx, cacheKey, jsonData, time.Minute*10).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, person)
	}
}

func AddPerson(db db.Database, rdb redisdb.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var person Person
		err := c.ShouldBindJSON(&person)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Проверяем, существует ли пользователь с такими данными
		existingUser := PersonModel{}
		err = db.Get(&existingUser, "SELECT * FROM people WHERE user_name = $1 AND surname = $2 AND patronymic = $3 AND age = $4 AND gender = $5 AND nationality = $6", person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
		if err == nil {
			logger.Warn("Пользователь с такими данными уже существует", err.Error())
			c.JSON(400, gin.H{"error": err})
			return
		}

		row := db.QueryRow("INSERT INTO people (user_name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)

		var generatedID uint
		err = row.Scan(&generatedID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		person.ID = generatedID

		// Очищаем кэш для данного пользователя
		ctx := context.Background()
		cacheKey := "person_" + strconv.FormatUint(uint64(generatedID), 10)
		err = rdb.Del(ctx, cacheKey).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"id":          generatedID,
			"user_name":   person.UserName,
			"surname":     person.Surname,
			"patronymic":  person.Patronymic,
			"age":         person.Age,
			"gender":      person.Gender,
			"nationality": person.Nationality,
		})
	}
}

func DeletePerson(db db.Database, rdb redisdb.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM people WHERE id = $1", id)
		if err != nil {
			logger.Error("Ошибка при удалении пользователя из базы: ", err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Очищаем кэш для данного пользователя
		ctx := context.Background()
		cacheKey := "person_" + id
		err = rdb.Del(ctx, cacheKey).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(200)
	}
}

func UpdatePerson(db db.Database, rdb redisdb.RedisClient) gin.HandlerFunc {
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
			logger.Error("Ошибка при обновлении пользователя в базе: ", err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Очищаем кэш для данного пользователя
		ctx := context.Background()
		cacheKey := "person_" + id
		err = rdb.Del(ctx, cacheKey).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"id":          id,
			"user_name":   person.UserName,
			"surname":     person.Surname,
			"patronymic":  person.Patronymic,
			"age":         person.Age,
			"gender":      person.Gender,
			"nationality": person.Nationality,
		})
	}
}
