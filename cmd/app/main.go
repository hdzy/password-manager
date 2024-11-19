package main

import (
	"password-manager/pkg/console"
	"time"
)

func main() {

	console.Start()

	<-time.NewTimer(30 * time.Second).C
}
