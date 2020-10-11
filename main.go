package main

import (
	"github.com/Phill93/DoorManager/alarm"
	"github.com/Phill93/DoorManager/aux"
	"github.com/Phill93/DoorManager/code"
	"github.com/Phill93/DoorManager/gate"
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
	g1 := gate.Gate{
		Name: "Tor 1",
	}
	a1 := aux.Aux{
		Number: 1,
	}

	go a.Watch()
	go wiegand.InitReader(&pad)
	go code2.ListenForKey(&pad)
	go handler.Handler(&code2, &a)
	go g1.ListenForOpen(&handler)
	go a1.ListenForTrigger(&handler)
	for {
		time.Sleep(1 * time.Second)
	}
}
