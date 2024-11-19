package main

import (
	"github.com/nsf/termbox-go"
	"strconv"
	"time"
)

func main() {
	value := 100
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputEsc)

	termbox.SetCell(0, 0, rune(value), 0, termbox.ColorBlue)

	//loop:
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			value++
			strValue := strconv.Itoa(value)
			//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//fmt.Println(byteValue)
			for i := 0; i < len(strValue); i++ {
				termbox.SetCell(i, 0, rune(strValue[i]), 0, termbox.ColorBlue)
			}
			termbox.Flush()
		}
	}

	<-time.NewTimer(1000 * time.Second).C
}
