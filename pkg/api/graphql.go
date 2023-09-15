package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/jmoiron/sqlx"
)

type PersonModel struct {
	ID          uint   `db:"id" json:"id"`
	UserName    string `db:"user_name" json:"user_name"`
	Surname     string `db:"surname" json:"surname"`
	Patronymic  string `db:"patronymic" json:"patronymic"`
	Age         int    `db:"age" json:"age"`
	Gender      string `db:"gender" json:"gender"`
	Nationality string `db:"nationality" json:"nationality"`
}

var db *sqlx.DB

// Определение типа PersonModel
var PersonType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PersonModel",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user_name": &graphql.Field{
				Type: graphql.String,
			},
			"surname": &graphql.Field{
				Type: graphql.String,
			},
			"patronymic": &graphql.Field{
				Type: graphql.String,
			},
			"age": &graphql.Field{
				Type: graphql.Int,
			},
			"gender": &graphql.Field{
				Type: graphql.String,
			},
			"nationality": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Определение GraphQL схемы
var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// Получение человека по ID
			"GetPerson": &graphql.Field{
				Type:        PersonType,
				Description: "Получить человека по ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Получаем значение ID из запроса
					id, ok := params.Args["id"].(int)
					if !ok {
						return nil, fmt.Errorf("Неверное значение ID")
					}

					db, ok := params.Context.Value("db").(*sqlx.DB)
					if !ok {
						return nil, fmt.Errorf("Не удалось получить доступ к базе данных")
					}

					query := "SELECT * FROM people WHERE id = $1"
					person := PersonModel{}
					err := db.Get(&person, query, id)
					if err != nil {
						return nil, err
					}

					return person, nil
				},
			},
			// Получение всех пользователей
			"AllPeople": &graphql.Field{
				Type:        graphql.NewList(PersonType),
				Description: "Получить всех пользователей",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db, ok := params.Context.Value("db").(*sqlx.DB)
					if !ok {
						return nil, fmt.Errorf("Не удалось получить доступ к базе данных")
					}

					var allPeople []PersonModel
					err := db.Select(&allPeople, "SELECT * FROM people")
					if err != nil {
						return nil, err
					}

					return allPeople, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			// Создать нового пользователя
			// Создать нового пользователя
			"CreatePerson": &graphql.Field{
				Type:        PersonType,
				Description: "Создать нового пользователя",
				Args: graphql.FieldConfigArgument{
					"user_name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"surname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"patronymic": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"gender": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"nationality": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["user_name"].(string)
					surname, _ := params.Args["surname"].(string)
					patronymic, _ := params.Args["patronymic"].(string)
					age, _ := params.Args["age"].(int)
					gender, _ := params.Args["gender"].(string)
					nationality, _ := params.Args["nationality"].(string)

					db, ok := params.Context.Value("db").(*sqlx.DB)
					if !ok {
						return nil, fmt.Errorf("Не удалось получить доступ к базе данных")
					}

					// Проверяем, существует ли пользователь с такими данными
					existingUser := PersonModel{}
					err := db.Get(&existingUser, "SELECT * FROM people WHERE user_name = $1 AND surname = $2 AND patronymic = $3 AND age = $4 AND gender = $5 AND nationality = $6", name, surname, patronymic, age, gender, nationality)
					if err == nil {
						return existingUser, nil
					}

					newPerson := PersonModel{
						UserName:    name,
						Surname:     surname,
						Patronymic:  patronymic,
						Age:         age,
						Gender:      gender,
						Nationality: nationality,
					}

					row := db.QueryRow(`INSERT INTO people (user_name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, newPerson.UserName, newPerson.Surname, newPerson.Patronymic, newPerson.Age, newPerson.Gender, newPerson.Nationality)

					var generatedID uint
					err = row.Scan(&generatedID)
					if err != nil {
						return nil, err
					}

					newPerson.ID = generatedID

					return newPerson, nil
				},
			},
			// Обновить данные пользователя
			"UpdatePerson": &graphql.Field{
				Type:        PersonType,
				Description: "Обновить данные пользователя",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"user_name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"surname": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"patronymic": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"age": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"gender": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"nationality": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					db, ok := params.Context.Value("db").(*sqlx.DB)
					if !ok {
						return nil, fmt.Errorf("Не удалось получить доступ к базе данных")
					}

					name, nameOk := params.Args["user_name"].(string)
					surname, surnameOk := params.Args["surname"].(string)
					patronymic, patronymicOk := params.Args["patronymic"].(string)
					age, ageOk := params.Args["age"].(int)
					gender, genderOk := params.Args["gender"].(string)
					nationality, nationalityOk := params.Args["nationality"].(string)

					updateQuery := "UPDATE people SET "
					values := []interface{}{id}

					if nameOk {
						updateQuery += "user_name = ?, "
						values = append(values, name)
					}
					if surnameOk {
						updateQuery += "surname = ?, "
						values = append(values, surname)
					}
					if patronymicOk {
						updateQuery += "patronymic = ?, "
						values = append(values, patronymic)
					}
					if ageOk {
						updateQuery += "age = ?, "
						values = append(values, age)
					}
					if genderOk {
						updateQuery += "gender = ?, "
						values = append(values, gender)
					}
					if nationalityOk {
						updateQuery += "nationality = ?, "
						values = append(values, nationality)
					}

					// Убираем последнюю запятую и пробел
					updateQuery = updateQuery[:len(updateQuery)-2]

					updateQuery += " WHERE id = ? RETURNING *"

					var updatedPerson PersonModel
					err := db.Get(&updatedPerson, updateQuery, values...)
					if err != nil {
						return nil, err
					}

					return updatedPerson, nil
				},
			},
			// Удалить пользователя
			"DeletePerson": &graphql.Field{
				Type:        PersonType,
				Description: "Удалить пользователя",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					db, ok := params.Context.Value("db").(*sqlx.DB)
					if !ok {
						return nil, fmt.Errorf("Не удалось получить доступ к базе данных")
					}

					deleteQuery := "DELETE FROM people WHERE id = $1 RETURNING *"

					var deletedPerson PersonModel
					err := db.Get(&deletedPerson, deleteQuery, id)
					if err != nil {
						return nil, err
					}

					return deletedPerson, nil
				},
			},
		},
	},
)

// Функция для обработки запросов GraphQL
func HandleGraphQL(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			Query string `json:"query"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := graphql.Params{
			Schema:         schema,
			RequestString:  requestBody.Query,
			RootObject:     map[string]interface{}{},
			VariableValues: map[string]interface{}{},
			OperationName:  "",
			Context:        c,
		}

		result := graphql.Do(params)
		if len(result.Errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": result.Errors})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
