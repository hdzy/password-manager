package console

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"os"
	"password-manager/pkg/policy"
	"reflect"
	"strings"
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

	if pickedValue == -1 {
		return -1
	}

	pickedPolicy := policies[pickedValue]

	if editPolicyMenu(pickedPolicy) == -1 {
		return showAllPolicies()
	}

	return 0
}

func editPolicyMenu(p *policy.Policy) int {
	ClearTerminal()
	RefreshTerminal()
	printItem(p.Name, termbox.ColorGreen, termbox.ColorBlack, 0)

	t := reflect.TypeOf(p).Elem()
	v := reflect.ValueOf(p).Elem()
	i := 0

	for ; i < t.NumField(); i++ {
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

		printItem(fmt.Sprintf("%s %v", tag, value), termbox.ColorWhite, termbox.ColorDefault, (i+1)*2)
	}

	RefreshTerminal()

	res, err := printButtons((i+2)*2+3, "Изменить", "Удалить")

	if err != nil {
		editPolicyMenu(p)
	}

	switch res {
	case 0:
		if editPolicy(p) == -1 {
			return editPolicyMenu(p)
		}
	case 1:
		if removePolicy(p) == -1 {
			return editPolicyMenu(p)
		}
	case -1:
		return -1
	}

	return 0
}

func editPolicy(p *policy.Policy) int {
	ClearTerminal()
	RefreshTerminal()
	printItem("Выберите поле для редактирования:", termbox.ColorGreen, termbox.ColorBlack, 0)

	t := reflect.TypeOf(p).Elem()
	v := reflect.ValueOf(p).Elem()

	cursorX := 0

	newValue := ""

	isWrite := false

	maxIndex := t.NumField()
	minIndex := 0
	currentItem := 0

	updateCursor := func() {
		termbox.SetCursor(cursorX, currentItem*2+2)
	}

	render := func() {
		ClearTerminal()
		i := 0
		printItem(p.Name, termbox.ColorGreen, termbox.ColorBlack, 0)
		for ; i < t.NumField(); i++ {
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
			if i == currentItem && !isWrite {
				printItem(fmt.Sprintf("%s %v", tag, value), termbox.ColorBlack, termbox.ColorWhite, (i+1)*2)
			} else if i != currentItem {
				printItem(fmt.Sprintf("%s %v", tag, value), termbox.ColorWhite, termbox.ColorDefault, (i+1)*2)
			} else if i == currentItem && isWrite {
				printItem(fmt.Sprintf("%s %v", tag, newValue), termbox.ColorWhite, termbox.ColorDefault, (i+1)*2)
			}
		}
		RefreshTerminal()
	}

	eventHandlerWrite := func() bool {
		for {
			event := termbox.PollEvent()

			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyEsc:
					return false
				case termbox.KeyEnter:
					return true
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					if len(newValue) > 0 {
						newValue = newValue[:len(newValue)-1]
						cursorX--
						updateCursor()
					}
					render()
				default:
					if event.Ch != 0 {
						newValue += string(event.Ch)
						cursorX++
						updateCursor()
						render()
					}
				}
			}
		}
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
					field := t.Field(currentItem)
					cursorX = runewidth.StringWidth(field.Tag.Get("print")) + 3
					termbox.SetCursor(cursorX, currentItem*2+2)
					isWrite = true
					newValue = ""
					render()
					if eventHandlerWrite() {
						policies.UpdateByName(p.Name, policy.Fvm{field.Name: newValue})
					}
					termbox.HideCursor()
					isWrite = false
					render()
				}
			}
		}
	}

	ClearTerminal()
	RefreshTerminal()
	render()
	if eventHandler() {
		return 0
	} else {
		return -1
	}
	return 0 // Added return statement
}

func removePolicy(p *policy.Policy) int {
	ClearTerminal()
	RefreshTerminal()
	printItem(p.Name, termbox.ColorGreen, termbox.ColorBlack, 0)

	t := reflect.TypeOf(p).Elem()
	v := reflect.ValueOf(p).Elem()
	i := 0

	for ; i < t.NumField(); i++ {
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

		printItem(fmt.Sprintf("%s %v", tag, value), termbox.ColorWhite, termbox.ColorDefault, (i+1)*2)
	}

	RefreshTerminal()

	res, _ := printButtons((i+2)*2+3, "Изменить", "Удалить")

	switch res {
	case 0:
	case 1:
	case -1:
		editPolicy(p)
	}

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
	x := 2

	for _, ch := range text {
		termbox.SetCell(x, y, ch, fg, bg)
		x += runewidth.RuneWidth(ch)
	}
}

func printButtons(y int, buttons ...string) (int, error) {
	width, _ := termbox.Size()

	totalButtonWidth := runewidth.StringWidth(strings.Join(buttons, "")) + len(buttons)*4
	spaceBetween := (width - totalButtonWidth - 4) / (len(buttons) - 1)

	currentItem := 0

	maxIndex := len(buttons) - 1
	minIndex := 0

	render := func() {
		x := 2
		for i, btn := range buttons {
			btn = "[ " + btn + " ]"
			for _, ch := range btn {
				if i == currentItem {
					termbox.SetCell(x, y, ch, termbox.ColorBlack, termbox.ColorWhite)
				} else {
					termbox.SetCell(x, y, ch, termbox.ColorWhite, termbox.ColorDefault)
				}
				x += runewidth.RuneWidth(ch)
			}
			x += spaceBetween
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
				case termbox.KeyArrowRight:
					if currentItem < maxIndex {
						currentItem++
						render()
					}
				case termbox.KeyArrowLeft:
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

	render()

	if eventHandler() {
		return currentItem, nil
	} else {
		return -1, nil
	}
}
