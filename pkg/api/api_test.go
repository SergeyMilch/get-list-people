package api

import (
	"testing"
)

func TestGetAge(t *testing.T) {
	age, err := GetAge("Dmitriy")
	if err != nil {
		t.Errorf("Ошибка при получении возраста: %s", err)
	}

	if age < 0 || age > 120 {
		t.Errorf("Ожидался корректный возраст, получено: %d", age)
	}
}

func TestGetGender(t *testing.T) {
	gender, err := GetGender("Dmitriy")
	if err != nil {
		t.Errorf("Ошибка при получении пола: %s", err)
	}

	if gender != "male" && gender != "female" {
		t.Errorf("Ожидался корректный пол, получено: %s", gender)
	}
}

func TestGetNationality(t *testing.T) {
	nationality, err := GetNationality("Dmitriy")
	if err != nil {
		t.Errorf("Ошибка при получении национальности: %s", err)
	}

	if nationality == "" {
		t.Errorf("Ожидалась корректная национальность, получено: %s", nationality)
	}
}
