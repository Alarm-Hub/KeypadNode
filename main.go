package main

import (
	"github.com/Phill93/DoorManager/code"
	//"github.com/Phill93/DoorManager/communicator"
	//"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/Phill93/DoorManager/wiegand"
	"time"
)

func main() {
	log.Infof("Application started at %s", time.Now())

	//cfg := config.Config()

	//c := communicator.NewCommunicator(cfg.GetString("access_token"), cfg.GetString("refresh_token"), cfg.GetString("controller_url"))
	//err := c.VerifyAccess()
	//if err != nil {
	//  panic(err)
	//}

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
			log.Debugf("Received code %s", <-codeCode)
		}
	}(&code2)

	go wiegand.InitReader(&pad)

	for {
		time.Sleep(1 * time.Second)
	}
}
