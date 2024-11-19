package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

func main() {
	value := 10
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputEsc)

	termbox.SetCell(0, 0, rune(value), 0, termbox.ColorBlue)
	termbox.Flush()
	//termbox.PollEvent()

	//loop:
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			value++
			termbox.SetCell(0, 0, rune(value), 0, termbox.ColorBlue)
			termbox.Sync()
			termbox.Flush()
		}
	}

	<-time.NewTimer(1000 * time.Second).C
}
