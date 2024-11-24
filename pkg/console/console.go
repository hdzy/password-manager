package console

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
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

	menuItems := []string{"Пункт 1", "Пункт 2", "Пункт 3", "Выход"}
	currentItem := 0

	ClearTerminal()
	pickedValue := initMenu(menuItems, currentItem)

	ClearTerminal()
	RefreshTerminal()
	fmt.Println("Выбранный пункт:", pickedValue+1)

	<-time.NewTimer(10 * time.Second).C
}

func initMenu(menuItems []string, currentItem int) int {
	maxIndex := len(menuItems) - 1
	minIndex := 0

	render := func() {
		for i, item := range menuItems {
			y := i + 1
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
	for x, ch := range text {
		termbox.SetCell(x/2, y, ch, fg, bg)
	}
}
