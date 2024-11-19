package console

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"password-manager/iternal/storage"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Очистка консоли
func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func NewPrompt(title string, items ...string) (string, error) {
	prompt := promptui.Select{
		Label:    title,
		HideHelp: true,
		Items:    items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return "", err
	}

	return result, nil
}

func Start() {
	expiredPasswords, err := storage.GetExpiredPasswords()
	if err != nil {
		fmt.Printf("Ошибка при проверке просроченных паролей: %v\n", err)
	} else if len(expiredPasswords) > 0 {
		fmt.Printf("\n⚠️  У вас есть %d просроченных паролей! Проверьте их в меню \"Просмотреть просроченные пароли\".\n\n", len(expiredPasswords))
	}

	for {
		startItems := []string{
			"Добавить новый пароль",
			"Посмотреть сохраненные пароли",
			"Найти пароль",
			"Просмотреть просроченные пароли",
			"Выход",
		}
		result, err := NewPrompt("Выберите действие:\n", startItems...)
		if err != nil {
			panic(err)
		}

		clearConsole()

		switch result {
		case startItems[0]:
			NewPass()
		case startItems[1]:
			ShowAllPasswords()
		case startItems[2]:
			FindPassword()
		case startItems[3]:
			ShowExpiredPasswords()
		case startItems[4]:
			fmt.Println("Выход из программы...")
			return
		}
	}
}

func NewPass() {
	reader := bufio.NewReader(os.Stdin)

	readRequiredString := func(prompt string) string {
		for {
			fmt.Print(prompt)
			input, _ := reader.ReadString('\n')
			trimmed := strings.TrimSpace(input)
			if trimmed != "" {
				return trimmed
			}
			fmt.Println("Поле обязательно для заполнения. Пожалуйста, введите значение.")
		}
	}

	resource := readRequiredString("Введите ресурс: ")
	identifier := readRequiredString("Введите идентификатор: ")
	password := readRequiredString("Введите пароль: ")

	var remindDate *string
	var remindDays *int
	for {
		daysStr := readRequiredString("Через сколько дней сменить пароль: ")

		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 {
			fmt.Println("Введите корректное количество дней (целое положительное число).")
			continue
		}

		date := time.Now().AddDate(0, 0, days)
		formattedDate := date.Format("2006-01-02")
		remindDate = &formattedDate
		remindDays = &days
		break
	}

	data := storage.PasswordData{
		Resource:   resource,
		Identifier: identifier,
		Password:   password,
		RemindDate: remindDate,
		RemindDays: remindDays,
		CreatedAt:  time.Now(),
	}

	err := storage.SavePassword(data)
	if err != nil {
		fmt.Printf("Ошибка сохранения данных: %v\n", err)
		return
	}

	fmt.Println("\nВаши данные сохранены.")
	WaitForEnter()
}

func ShowAllPasswords() {
	passwords, err := storage.GetAllPasswords()
	if err != nil {
		fmt.Printf("Ошибка получения паролей: %v\n", err)
		return
	}

	if len(passwords) == 0 {
		fmt.Println("Нет сохраненных паролей.")
	} else {
		fmt.Println("\nВсе сохраненные пароли:\n\n")
		for _, p := range passwords {
			fmt.Printf("Ресурс: %s, Идентификатор: %s, Пароль: %s\n", p.Resource, p.Identifier, p.Password)
			if p.RemindDate != nil {
				fmt.Println("Напоминание о смене пароля:", *p.RemindDate)
			}
			fmt.Println()
		}
	}
	WaitForEnter()
}

func FindPassword() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите название ресурса или идентификатор для поиска: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	passwords, err := storage.FindPassword(query)
	if err != nil {
		fmt.Printf("Ошибка поиска пароля: %v\n", err)
		return
	}

	if len(passwords) == 0 {
		fmt.Println("Пароль не найден.")
	} else {
		fmt.Println("\nНайденные пароли:\n")
		for _, p := range passwords {
			fmt.Printf("Ресурс: %s, Идентификатор: %s, Пароль: %s\n", p.Resource, p.Identifier, p.Password)
			if p.RemindDate != nil {
				fmt.Println("Напоминание о смене пароля:", *p.RemindDate)
			}
			fmt.Println()
		}
	}
	WaitForEnter()
}

func WaitForEnter() {
	fmt.Println("\nНажмите Enter для возврата в главное меню...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	clearConsole()
}

func ShowExpiredPasswords() {
	passwords, err := storage.GetExpiredPasswords()
	if err != nil {
		fmt.Printf("Ошибка получения данных: %v\n", err)
		return
	}

	if len(passwords) == 0 {
		fmt.Println("Нет просроченных паролей.")
	} else {
		for {
			fmt.Println("\nПросроченные пароли:\n")
			items := make([]string, len(passwords)+1)
			for i, p := range passwords {
				items[i] = fmt.Sprintf("Ресурс: %s, Идентификатор: %s", p.Resource, p.Identifier)
			}
			items[len(passwords)] = "Назад"

			selected, err := NewPrompt("Изменить пароль:\n", items...)
			if err != nil {
				fmt.Printf("Ошибка выбора: %v\n", err)
				return
			}

			if selected == "Назад" {
				return
			}

			selectedIndex := -1
			for i, item := range items {
				if item == selected {
					selectedIndex = i
					break
				}
			}

			if selectedIndex >= 0 && selectedIndex < len(passwords) {
				editPassword(&passwords[selectedIndex])
			}
		}
	}
	WaitForEnter()
}

func editPassword(password *storage.PasswordData) {
	for {
		fmt.Printf("\nРедактирование пароля:\nРесурс: %s\nИдентификатор: %s\nПароль: %s\n",
			password.Resource, password.Identifier, password.Password)
		if password.RemindDate != nil {
			fmt.Println("Дата напоминания:", *password.RemindDate)
		}

		options := []string{
			"Изменить ресурс",
			"Изменить идентификатор",
			"Изменить пароль",
			"Назад",
		}

		choice, err := NewPrompt("Выберите действие\n", options...)
		if err != nil {
			fmt.Printf("Ошибка выбора: %v\n", err)
			return
		}

		switch choice {
		case "Изменить ресурс":
			fmt.Print("Введите новый ресурс: ")
			newResource := readInput()
			password.Resource = newResource
		case "Изменить идентификатор":
			fmt.Print("Введите новый идентификатор: ")
			newIdentifier := readInput()
			password.Identifier = newIdentifier
		case "Изменить пароль":
			fmt.Print("Введите новый пароль: ")
			newPassword := readInput()
			password.Password = newPassword
			if password.RemindDays != nil {
				newDate := time.Now().AddDate(0, 0, *password.RemindDays).Format("2006-01-02")
				password.RemindDate = &newDate
			}
		case "Назад":
			err := storage.UpdatePassword(*password)
			if err != nil {
				fmt.Printf("Ошибка сохранения изменений: %v\n", err)
			} else {
				fmt.Println("Изменения успешно сохранены.")
			}
			return
		}
	}
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
