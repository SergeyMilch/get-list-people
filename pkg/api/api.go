package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SergeyMilch/get-list-people-effective-mobile/pkg/logger"
)

func GetAge(name string) (uint8, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Ошибка при запросе возраста:", err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Error("Ошибка при разборе JSON ответа:", err.Error())
		return 0, err
	}

	age := uint8(data["age"].(float64))
	return age, nil
}

func GetGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Ошибка при запросе пола:", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Error("Ошибка при разборе JSON ответа:", err.Error())
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
		logger.Error("Ошибка при запросе национальности:", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Error("Ошибка при разборе JSON ответа:", err.Error())
		return "", err
	}

	countries, ok := data["country"].([]interface{})
	if !ok || len(countries) == 0 {
		logger.Error("Не удалось получить данные о странах", err.Error())
		return "", err
	}

	var mostProbableCountryID string
	var maxProbability float64 = 0

	for _, country := range countries {
		countryData, ok := country.(map[string]interface{})
		if !ok {
			logger.Error("Не удалось получить данные о стране", err.Error())
			return "", err
		}

		probability, ok := countryData["probability"].(float64)
		if !ok {
			logger.Error("Не удалось получить вероятность", err.Error())
			return "", err
		}

		if probability > maxProbability {
			maxProbability = probability
			mostProbableCountryID, ok = countryData["country_id"].(string)
			if !ok {
				logger.Error("Не удалось получить ID страны", err.Error())
				return "", err
			}
		}
	}

	if mostProbableCountryID == "" {
		logger.Error("Не удалось найти наиболее вероятную страну", err.Error())
		return "", err
	}

	return mostProbableCountryID, nil
}
