package main

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

func notify() {

	beeep.AppName = "My App Name"
	title := "Title"
	message := "Message body"

	icon := "icons/warning.png"

	fmt.Println("This")

	// err := beeep.Notify(title, message, icon)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := beeep.Alert(title, message, icon)
	if err != nil {
		panic(err)
	}
}
