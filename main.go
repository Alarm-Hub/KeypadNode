package main

import (
	"github.com/Phill93/DoorManager/alarm"
	"github.com/Phill93/DoorManager/code"
	"github.com/Phill93/DoorManager/log"
	"github.com/Phill93/DoorManager/wiegand"
	"time"
)

func main() {
	log.Infof("Application started at %s", time.Now())
	pad := wiegand.Keypad{}
	code2 := code.Code{}
	handler := code.Handler{}
	a := alarm.Alarm{}

	go a.Watch()
	go wiegand.InitReader(&pad)
	go code2.ListenForKey(&pad)
	handler.Handler(&code2, &a)
	for {
		time.Sleep(1 * time.Second)
	}
}
