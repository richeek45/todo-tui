package main

import (
	"log"

	"github.com/gen2brain/beeep"
)

func notify() {

	beeep.AppName = "My App Name"
	title := "Low Battery"
	message := "Your Mac will sleep soon unless plugged into a power outlet"

	icon := "icons/low_battery.png"

	err := beeep.Notify(title, message, icon)
	if err != nil {
		log.Fatal(err)
	}

	// err := beeep.Alert(title, message, icon)
	// if err != nil {
	// 	panic(err)
	// }
}
