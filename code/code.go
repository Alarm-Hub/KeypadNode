package code

import (
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/Phill93/DoorManager/wiegand"
	"strings"
	"time"
)

type Code struct {
	digits    []string
	startTime time.Time
	listeners map[string][]chan string
}

func (c *Code) AddListener(e string, ch chan string) {
	if c.listeners == nil {
		c.listeners = make(map[string][]chan string)
	}
	log.Debugf("Adding channel for event %s", e)
	if _, ok := c.listeners[e]; ok {
		c.listeners[e] = append(c.listeners[e], ch)
	} else {
		c.listeners[e] = []chan string{ch}
	}
}

func (c *Code) RemoveListener(e string, ch chan string) {
	log.Debugf("Removing channel for event %s", e)
	if _, ok := c.listeners[e]; ok {
		for i := range c.listeners[e] {
			if c.listeners[e][i] == ch {
				c.listeners[e] = append(c.listeners[e][:i], c.listeners[e][i+1:]...)
				break
			}
		}
	}
}

func (c *Code) Emit(e string, response string) {
	log.Debugf("Event %s is emitted!", e)
	if _, ok := c.listeners[e]; ok {
		for _, handler := range c.listeners[e] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}

func (c *Code) Input(key string) {
	if c.startTime.IsZero() {
		c.startTime = time.Now()
		log.Debugf("Start time empty! Set to now %s", c.startTime)
	}
	cfg := config.Config()
	if checkTimeout(c.startTime, cfg.GetInt("pad_timeout")) {
		c.Clear()
	} else {
		c.startTime = time.Now()
	}

	if key == "ENT" {
		log.Debug("Enter Key received!")
		c.Emit("code", strings.Join(c.digits, ""))
		c.Clear()
	} else if key == "ESC" {
		log.Debug("Escape Key received!")
		c.Clear()
	} else {
		log.Debugf("Received key %s!", key)
		c.digits = append(c.digits, key)
	}
}

func (c *Code) Clear() {
	log.Info("Clear code!")
	c.digits = nil
	c.startTime = time.Time{}
}

func checkTimeout(timestamp time.Time, timeout int) bool {
	now := time.Now()
	offset := now.Sub(timestamp)
	if offset.Seconds() > float64(timeout) {
		log.Infof("Timeout reached! %f > %f", offset.Seconds(), float64(timeout))
		return true
	} else {
		log.Infof("Timeout not reached! %f !> %f", offset.Seconds(), float64(timeout))
		return false
	}
}

func (c *Code) ListenForKey(p *wiegand.Keypad) {
	events := make(chan string)
	p.AddListener("key", events)
	for {
		log.Debugf("Code: Waiting for key")
		event := <-events
		log.Debugf("Code: Received key %s", event)
		c.Input(event)
	}
}
