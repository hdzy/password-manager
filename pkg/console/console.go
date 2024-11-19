package console

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
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

func NewPrompt(items ...string) (string, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return "", err
	}

	return result, nil
}

func Start() {
	startItems := []string{"Добавить новый пароль", "Посмотреть сохраненные пароли", "Найти пароль"}
	result, err := NewPrompt(startItems...)
	if err != nil {
		panic(err)
	}

	// Очистка консоли после выбора
	clearConsole()

	switch result {
	case startItems[0]:
		NewPass()
	case startItems[1]:
		fmt.Println("all pass")
	case startItems[2]:
		fmt.Println("find pass")
	}
}

func NewPass() {
	reader := bufio.NewReader(os.Stdin)

	// Функция для запроса строки с проверкой заполненности
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

	// Функция для необязательного ввода
	readOptionalString := func(prompt string) string {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		return strings.TrimSpace(input)
	}

	// Запрашиваем обязательные данные
	resource := readRequiredString("Введите ресурс: ")
	identifier := readRequiredString("Введите идентификатор: ")
	password := readRequiredString("Введите пароль: ")

	var remindDate *time.Time // Указатель, чтобы обработать отсутствие даты
	for {
		daysStr := readOptionalString("Через сколько дней сменить пароль (Enter для отмены): ")
		if daysStr == "" {
			break
		}

		// Конвертация дней в дату
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 {
			fmt.Println("Введите корректное количество дней (целое положительное число).")
			continue
		}

		// Рассчитываем дату напоминания
		date := time.Now().AddDate(0, 0, days)
		remindDate = &date
		break
	}

	// Вывод результата
	fmt.Println("\nВаши данные:")
	fmt.Println("Ресурс:", resource)
	fmt.Println("Идентификатор:", identifier)
	fmt.Println("Пароль:", password)
	if remindDate != nil {
		fmt.Println("Напоминание о смене пароля:", remindDate.Format("2006-01-02"))
	} else {
		fmt.Println("Напоминание отключено.")
	}
}
