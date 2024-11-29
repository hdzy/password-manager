package console

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"os"
	"password-manager/pkg/policy"
	"reflect"
)

type (
	MenuAction int
)

const (
	prev MenuAction = iota
	next
)

var policies = policy.Policies{}

func Init() {
	if err := termbox.Init(); err != nil {
		fmt.Println("Ошибка инициализации termbox:", err)
		os.Exit(1)
	}
	defer termbox.Close()

	policies.Load()

	if startMenu() == -1 {
		return
	}

	termbox.PollEvent()
}

func startMenu() int {
	menuItems := []string{"Управление парольными политиками", "Управление паролями"}
	currentItem := 0

	switch pickedValue := initMenu(menuItems, currentItem, "Главное меню"); pickedValue {
	case 0:
		if policyMenu() == -1 {
			return startMenu()
		}
	case 1:
		if passwordMenu() == -1 {
			return startMenu()
		}
	case -1:
		return -1
	}

	ClearTerminal()
	RefreshTerminal()

	return 1
}

func policyMenu() int {
	menuItems := []string{"Посмотреть все парольные политики", "Создать новую парольную политику"}
	currentItem := 0

	switch pickedValue := initMenu(menuItems, currentItem, "Меню управления парольными политиками"); pickedValue {
	case -1:
		return -1
	case 0:
		if showAllPolicies() == -1 {
			return policyMenu()
		}
	}

	ClearTerminal()
	RefreshTerminal()

	return 1
}

func showAllPolicies() int {
	menuItems := make([]string, len(policies))

	for i, el := range policies {
		menuItems[i] = el.Name
	}

	pickedValue := initMenu(menuItems, 0, "Парольные политики (выберите для редактирования):")
	pickedPolicy := policies[pickedValue]

	if editPolicy(pickedPolicy) == -1 {
		return showAllPolicies()
	}

	return 0
}

func editPolicy(p *policy.Policy) int {
	ClearTerminal()
	RefreshTerminal()
	printItem(p.Name, termbox.ColorGreen, termbox.ColorBlack, 0)

	t := reflect.TypeOf(p).Elem()
	v := reflect.ValueOf(p).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("print")
		var value interface{}

		if v.Field(i).Kind() == reflect.Slice && v.Field(i).Type().Elem().Kind() == reflect.Int32 {
			runes := v.Field(i).Interface().([]rune)
			value = ""
			for _, ch := range runes {
				value = value.(string) + string(ch)
			}
		} else {
			value = v.Field(i).Interface()
		}

		printItem(fmt.Sprintf("%s %v", tag, value), termbox.ColorWhite, termbox.ColorDefault, i+2)
	}

	RefreshTerminal()

	// Wait for a key press before returning
	termbox.PollEvent()

	return 0
}

func passwordMenu() int {
	menuItems := []string{"Посмотреть все пароли", "Создать новый пароль"}
	currentItem := 0

	pickedValue := initMenu(menuItems, currentItem, "Управление паролями")

	if pickedValue == -1 {
		return -1
	}

	ClearTerminal()
	RefreshTerminal()

	fmt.Println("Выбранный пункт:", pickedValue+1)

	return 1
}

func initMenu(menuItems []string, currentItem int, title string) int {
	maxIndex := len(menuItems) - 1
	minIndex := 0

	render := func() {
		printItem(title, termbox.ColorGreen, termbox.ColorBlack, 0)
		for i, item := range menuItems {
			y := i + 2
			if i == currentItem {
				printItem(item, termbox.ColorBlack, termbox.ColorWhite, y)
			} else {
				printItem(item, termbox.ColorWhite, termbox.ColorDefault, y)
			}
		}
		RefreshTerminal()
	}

	eventHandler := func() bool {
		for {
			event := termbox.PollEvent()

			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyEsc:
					return false
				case termbox.KeyArrowDown:
					if currentItem < maxIndex {
						currentItem++
						render()
					}
				case termbox.KeyArrowUp:
					if currentItem > minIndex {
						currentItem--
						render()
					}
				case termbox.KeyEnter:
					return true
				}
			}
		}
	}

	ClearTerminal()
	RefreshTerminal()
	render()
	if eventHandler() {
		return currentItem
	} else {
		return -1
	}
}

func ClearTerminal() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func RefreshTerminal() {
	termbox.Flush()
}

// Функция для вывода строки целиком
func printItem(text string, fg, bg termbox.Attribute, y int) {
	x := 0

	for _, ch := range text {
		termbox.SetCell(x, y, ch, fg, bg)
		x += runewidth.RuneWidth(ch)
	}
}
