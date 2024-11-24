package console

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
	"time"
)

type (
	MenuAction int
)

const (
	prev MenuAction = iota
	next
)

func Init() {
	if err := termbox.Init(); err != nil {
		fmt.Println("Ошибка инициализации termbox:", err)
		os.Exit(1)
	}
	defer termbox.Close()

	if startMenu() == -1 {
		return
	}

	<-time.NewTimer(10 * time.Second).C
}

func startMenu() int {
	menuItems := []string{"Управление парольными политиками", "Управление паролями"}
	currentItem := 0

	switch pickedValue := initMenu(menuItems, currentItem); pickedValue {
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

	pickedValue := initMenu(menuItems, currentItem)

	if pickedValue == -1 {
		return -1
	}

	ClearTerminal()
	RefreshTerminal()

	fmt.Println("Выбранный пункт:", pickedValue+1)

	return 1
}

func passwordMenu() int {
	menuItems := []string{"Посмотреть все пароли", "Создать новый пароль"}
	currentItem := 0

	pickedValue := initMenu(menuItems, currentItem)

	if pickedValue == -1 {
		return -1
	}

	ClearTerminal()
	RefreshTerminal()

	fmt.Println("Выбранный пункт:", pickedValue+1)

	return 1
}

func initMenu(menuItems []string, currentItem int) int {
	maxIndex := len(menuItems) - 1
	minIndex := 0

	render := func() {
		for i, item := range menuItems {
			y := i * 2
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
	// Отображаем всю строку целиком
	text = strings.ReplaceAll(text, " ", "  ")
	for x, ch := range text {
		termbox.SetCell(x/2, y, ch, fg, bg)
	}
}
