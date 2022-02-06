package main

import (
  "github.com/Phill93/DoorManager/code"
  "github.com/Phill93/DoorManager/communicator"
  "github.com/Phill93/DoorManager/config"

  //"github.com/Phill93/DoorManager/communicator"
  //"github.com/Phill93/DoorManager/config"
  "github.com/Phill93/DoorManager/log"
  "github.com/Phill93/DoorManager/wiegand"
  "time"
)

func main() {
	log.Infof("Application started at %s", time.Now())

	cfg := config.Config()

	c := communicator.NewCommunicator(cfg.GetString("controller_url"))
	go func() {
	  for {
      ok, err := c.VerifyAccess()
      if err != nil {
        log.Error(err)
      }
      if !ok {
        err = c.Refresh()
        if err != nil {
          log.Error(err)
        }
      }
      time.Sleep(time.Second * 20)
    }
  }()

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

	go func(c *code.Code, com *communicator.Communicator) {
		codeCode := make(chan string)
		c.AddListener("ready", codeCode)
		for {
			log.Debugf("Waiting for code!")
			code := <-codeCode
			log.Debugf("Received code %s", code)
			ok, err := com.ValidateCode(code)
			if err != nil {
			  log.Error(err)
      }
      if !ok {
        log.Errorf("Code %s is invalid!", code)
      } else {
        log.Info("Code %s is valid!", code)
      }
		}
	}(&code2, &c)

	go wiegand.InitReader(&pad)

	for {
		time.Sleep(1 * time.Second)
	}
}
