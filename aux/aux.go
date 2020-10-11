package aux

import (
	"github.com/Phill93/DoorManager/code"
	"github.com/Phill93/DoorManager/log"
)

type Aux struct {
	Number int
}

func (a *Aux) ListenForOn(ch *code.Handler) {
	events := make(chan string)
	ch.AddListener("result", events)
	for {
		log.Debugf("Aux %d: Waiting for On Event", a.Number)
		event := <-events
		log.Debugf("Aux %d: Received On Event %s", a.Number, event)
		if event == "valid" {
			a.On()
		}
	}
}

func (a *Aux) ListenForOff(ch *code.Handler) {
	events := make(chan string)
	ch.AddListener("result", events)
	for {
		log.Debugf("Aux %d: Waiting for Off Event", a.Number)
		event := <-events
		log.Debugf("Aux %d: Received Aux Off %s", a.Number, event)
		if event == "valid" {
			a.Off()
		}
	}
}

func (a *Aux) ListenForTrigger(ch *code.Handler) {
	events := make(chan string)
	ch.AddListener("result", events)
	for {
		log.Debugf("Aux %d: Waiting for Trigger Event", a.Number)
		event := <-events
		log.Debugf("Aux %d: Received Trigger Event %s", a.Number, event)
		if event == "valid" {
			a.Trigger()
		}
	}
}
