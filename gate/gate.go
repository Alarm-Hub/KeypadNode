package gate

import (
	"github.com/Phill93/DoorManager/code"
	"github.com/Phill93/DoorManager/log"
)

type Gate struct {
	Name     string
	PinOpen  int
	PinClose int
	PinStop  int
}

func (g *Gate) ListenForOpen(ch *code.Handler) {
	events := make(chan string)
	ch.AddListener("result", events)
	for {
		log.Debugf("Gate %s: Waiting for Open Event", g.Name)
		event := <-events
		log.Debugf("Gate %s: Received Open Event %s", g.Name, event)
		if event == "valid" {
			g.open()
		}
	}
}

func (g *Gate) ListenForClose(ch *code.Handler) {
	events := make(chan string)
	ch.AddListener("result", events)
	for {
		log.Debugf("Gate %s: Waiting for Close Event", g.Name)
		event := <-events
		log.Debugf("Gate %s: Received Close Event %s", g.Name, event)
	}
}

func (g *Gate) ListenForStop(ch *code.Handler) {
	events := make(chan string)
	ch.AddListener("result", events)
	for {
		log.Debugf("Gate %s: Waiting for Stop Event", g.Name)
		event := <-events
		log.Debugf("Gate %s: Received Stop Event %s", g.Name, event)
	}
}
