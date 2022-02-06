package wiegand

import (
	"github.com/Phill93/DoorManager/log"
)

var low uint32 = 0x0
var high uint32 = 0x0
var carry uint32 = 0
var lowWd chan bool
var highWd chan bool
var end = false

type Keypad struct {
	listeners map[string][]chan string
}

func (k *Keypad) AddListener(e string, ch chan string) {
	if k.listeners == nil {
		k.listeners = make(map[string][]chan string)
	}
	log.Debugf("Adding channel for event %s", e)
	if _, ok := k.listeners[e]; ok {
		k.listeners[e] = append(k.listeners[e], ch)
	} else {
		k.listeners[e] = []chan string{ch}
	}
}

func (k *Keypad) RemoveListener(e string, ch chan string) {
	log.Debugf("Removing channel for event %s", e)
	if _, ok := k.listeners[e]; ok {
		for i := range k.listeners[e] {
			if k.listeners[e][i] == ch {
				k.listeners[e] = append(k.listeners[e][:i], k.listeners[e][i+1:]...)
				break
			}
		}
	}
}

func (k *Keypad) Emit(e string, response string) {
	log.Debugf("Event %s is emitted with Value %s!", e, response)
	if _, ok := k.listeners[e]; ok {
		for _, handler := range k.listeners[e] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}
