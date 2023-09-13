package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAge(name string) (int, error) {
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

	age := int(data["age"].(int))
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
