package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SergeyMilch/get-list-people-effective-mobile/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetPeople(t *testing.T) {
	// Создаем контроллер Gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для Database
	mockDB := mocks.NewMockDatabase(ctrl)

	// Создаем мок для RedisClient
	mockRedis := mocks.NewMockRedisClient(ctrl)

	// Создаем тестовый контекст и запрос
	req := httptest.NewRequest(http.MethodGet, "/people", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Устанавливаем ожидаемое значение в Redis mock
	expectedPeople := []Person{{ID: 1, UserName: "JohnDoe", Surname: "Doe", Patronymic: "Jr.", Age: 30, Gender: "Male", Nationality: "US"}}
	expectedJsonData, _ := json.Marshal(expectedPeople)
	mockRedis.EXPECT().Get(gomock.Any(), "all_people").Return(redis.NewStringResult(string(expectedJsonData), nil))

	// Помещаем моки в контекст Gin
	ctx.Set("db", mockDB)
	ctx.Set("rdb", mockRedis)

	// Устанавливаем контекст для запроса
	req = req.WithContext(ctx)

	// Вызываем функцию GetPeople
	GetPeople(mockDB, mockRedis)(ctx)

	// Проверяем, что код состояния HTTP равен http.StatusOK
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что возвращенный JSON имеет ожидаемую структуру и содержимое
	var response []Person
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPeople, response)
}

func TestGetPersonByID(t *testing.T) {
	// Создаем контроллер Gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для Database
	mockDB := mocks.NewMockDatabase(ctrl)

	// Создаем мок для RedisClient
	mockRedis := mocks.NewMockRedisClient(ctrl)

	// Создаем тестовый контекст и запрос
	req := httptest.NewRequest(http.MethodGet, "/people/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Устанавливаем ожидаемое значение в Redis mock
	expectedPerson := Person{ID: 1, UserName: "JohnDoe", Surname: "Doe", Patronymic: "Jr.", Age: 30, Gender: "Male", Nationality: "US"}
	expectedJsonData, _ := json.Marshal(expectedPerson)
	mockRedis.EXPECT().Get(gomock.Any(), "person_").Return(redis.NewStringResult(string(expectedJsonData), nil))

	// Помещаем моки в контекст Gin
	ctx.Set("db", mockDB)
	ctx.Set("rdb", mockRedis)

	// Устанавливаем контекст для запроса
	req = req.WithContext(ctx)

	// Вызываем функцию GetPersonByID
	GetPersonByID(mockDB, mockRedis)(ctx)

	// Проверяем, что код состояния HTTP равен http.StatusOK
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что возвращенный JSON имеет ожидаемую структуру и содержимое
	var response Person
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPerson, response)
}
