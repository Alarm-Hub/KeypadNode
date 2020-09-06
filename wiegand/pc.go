// +build !rpi

package wiegand

import (
	"github.com/eiannone/keyboard"
	"os"
)

func InitReader(pad *Keypad) {
	keyEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	for {
		event := <-keyEvents
		if event.Err != nil {
			panic(event.Err)
		}

		switch string(event.Rune) {
		case "1":
			pad.Emit("key", "1")
		case "2":
			pad.Emit("key", "2")
		case "3":
			pad.Emit("key", "3")
		case "4":
			pad.Emit("key", "4")
		case "5":
			pad.Emit("key", "5")
		case "6":
			pad.Emit("key", "6")
		case "7":
			pad.Emit("key", "7")
		case "8":
			pad.Emit("key", "8")
		case "9":
			pad.Emit("key", "9")
		case "0":
			pad.Emit("key", "0")
		}

		switch event.Key {
		case keyboard.KeyEsc:
			pad.Emit("key", "ESC")
		case keyboard.KeyEnter:
			pad.Emit("key", "ENT")
		case keyboard.KeyF12:
			os.Exit(0)
		}
	}
}
