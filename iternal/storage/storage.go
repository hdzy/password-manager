package storage

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type PasswordData struct {
	Resource   string    `json:"resource"`
	Identifier string    `json:"identifier"`
	Password   string    `json:"password"`
	RemindDate *string   `json:"remind_date,omitempty"`
	RemindDays *int      `json:"remind_days,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

var filePath = "./files/passwords.json"

// Сохранить пароль в файл
func SavePassword(data PasswordData) error {
	passwords, err := GetAllPasswords()
	if err != nil {
		return err
	}

	passwords = append(passwords, data)
	return saveAllPasswords(passwords)
}

// Сохранить все пароли
func saveAllPasswords(passwords []PasswordData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(passwords)
}

// Получить все пароли
func GetAllPasswords() ([]PasswordData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []PasswordData{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var passwords []PasswordData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&passwords)
	if err != nil {
		return nil, err
	}
	return passwords, nil
}

// Найти пароль по запросу
func FindPassword(query string) ([]PasswordData, error) {
	passwords, err := GetAllPasswords()
	if err != nil {
		return nil, err
	}

	var result []PasswordData
	for _, p := range passwords {
		if strings.Contains(p.Resource, query) || strings.Contains(p.Identifier, query) {
			result = append(result, p)
		}
	}
	return result, nil
}

// Получить просроченные пароли
func GetExpiredPasswords() ([]PasswordData, error) {
	passwords, err := GetAllPasswords()
	if err != nil {
		return nil, err
	}

	var expired []PasswordData
	for _, p := range passwords {
		if p.RemindDays != nil && p.RemindDate != nil {
			remindTime, err := time.Parse("2006-01-02", *p.RemindDate)
			if err != nil {
				continue
			}
			if remindTime.Before(time.Now()) {
				expired = append(expired, p)
			}
		}
	}
	return expired, nil
}

// Обновить пароль
func UpdatePassword(updatedPassword PasswordData) error {
	passwords, err := GetAllPasswords()
	if err != nil {
		return err
	}

	// Найти пароль и обновить
	found := false
	for i, p := range passwords {
		if p.Resource == updatedPassword.Resource && p.Identifier == updatedPassword.Identifier {
			passwords[i] = updatedPassword
			if updatedPassword.RemindDays != nil {
				newDate := time.Now().AddDate(0, 0, *updatedPassword.RemindDays).Format("2006-01-02")
				passwords[i].RemindDate = &newDate
			}
			found = true
			break
		}
	}

	if !found {
		return errors.New("запись не найдена")
	}

	// Перезапись данных в файл
	return saveAllPasswords(passwords)
}
