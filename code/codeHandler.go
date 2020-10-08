package code

import (
	"github.com/Phill93/DoorManager/alarm"
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"time"
)

type Handler struct {
	listeners map[string][]chan string
	lastAlarm time.Time
}

func (h *Handler) Handler(c *Code, a *alarm.Alarm) {
	go h.listenForCode(c)
	go h.listenForAlarm(a)
}

func (h *Handler) listenForCode(c *Code) {
	events := make(chan string)
	c.AddListener("code", events)
	for {
		log.Debugf("Code Handler: Waiting for code")
		event := <-events
		log.Debugf("Code Handler: Received code %s", event)
		h.checkCode(event)
	}
}

func (h *Handler) listenForAlarm(a *alarm.Alarm) {
	events := make(chan string)
	a.AddListener("alarm", events)
	for {
		log.Debugf("Code Handler: Waiting for alarm")
		event := <-events
		log.Debugf("Code Handler: Received alarm %s", event)
		h.lastAlarm = time.Now()
	}
}

func (h *Handler) checkCode(c string) {
	cfg := config.Config()
	if c == cfg.GetString("code") {
		log.Debug("Code Handler: Code is valid!")
		if !checkTimeout(h.lastAlarm, cfg.GetInt("timeout")) {
			log.Debug("Code Handler: Timeout not reached")
			h.Emit("result", "valid")
		} else {
			log.Debug("Code Handler: Timeout reached")
		}
	} else {
		log.Debug("Code Handler: Code is invalid!")
		h.Emit("result", "invalid")
	}
}

func (h *Handler) AddListener(e string, ch chan string) {
	if h.listeners == nil {
		h.listeners = make(map[string][]chan string)
	}
	log.Debugf("Adding channel for event %s", e)
	if _, ok := h.listeners[e]; ok {
		h.listeners[e] = append(h.listeners[e], ch)
	} else {
		h.listeners[e] = []chan string{ch}
	}
}

func (h *Handler) RemoveListener(e string, ch chan string) {
	log.Debugf("Removing channel for event %s", e)
	if _, ok := h.listeners[e]; ok {
		for i := range h.listeners[e] {
			if h.listeners[e][i] == ch {
				h.listeners[e] = append(h.listeners[e][:i], h.listeners[e][i+1:]...)
				break
			}
		}
	}
}

func (h *Handler) Emit(e string, response string) {
	log.Debugf("Event %s is emitted!", e)
	if _, ok := h.listeners[e]; ok {
		for _, handler := range h.listeners[e] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}
