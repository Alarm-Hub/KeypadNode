package alarm

import (
	"github.com/Phill93/DoorManager/log"
	"time"
)

type Alarm struct {
	lastAlarm time.Time
	listeners map[string][]chan string
}

func (a *Alarm) AddListener(e string, ch chan string) {
	if a.listeners == nil {
		a.listeners = make(map[string][]chan string)
	}
	log.Debugf("Adding channel for event %s", e)
	if _, ok := a.listeners[e]; ok {
		a.listeners[e] = append(a.listeners[e], ch)
	} else {
		a.listeners[e] = []chan string{ch}
	}
}

func (a *Alarm) RemoveListener(e string, ch chan string) {
	log.Debugf("Removing channel for event %s", e)
	if _, ok := a.listeners[e]; ok {
		for i := range a.listeners[e] {
			if a.listeners[e][i] == ch {
				a.listeners[e] = append(a.listeners[e][:i], a.listeners[e][i+1:]...)
				break
			}
		}
	}
}

func (a *Alarm) Emit(e string, response string) {
	log.Debugf("Event %s is emitted!", e)
	if _, ok := a.listeners[e]; ok {
		for _, handler := range a.listeners[e] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}
