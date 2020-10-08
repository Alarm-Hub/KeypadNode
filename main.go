package main

import (
	"github.com/Phill93/DoorManager/code"
	"github.com/Phill93/DoorManager/log"
	"github.com/Phill93/DoorManager/wiegand"
	"time"
)

func main() {
	log.Infof("Application started at %s", time.Now())
	pad := wiegand.Keypad{}
	code2 := code.Code{}

	go func(p *wiegand.Keypad, c *code.Code) {
		padKey := make(chan string)
		p.AddListener("key", padKey)
		for {
			select {
			case key := <-padKey:
				log.Printf("Got key %s", key)
				c.Input(key)
			}
		}
	}(&pad, &code2)

	go func(c *code.Code) {
		codeCode := make(chan string)
		c.AddListener("ready", codeCode)
		for {
			log.Debugf("Waiting for code!")
			code := <-codeCode
			log.Debugf("Received code %s", code)
		}
	}(&code2)

	go wiegand.InitReader(&pad)

	for {
		time.Sleep(1 * time.Second)
	}
}
