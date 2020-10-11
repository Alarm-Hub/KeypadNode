// +build rpi

package aux

import (
	"fmt"
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/warthog618/gpiod"
	"time"
)

func (a *Aux) On() {
	cfg := config.Config()
	log.Infof("Aux %d on", a.Number)
	c, _ := gpiod.NewChip("gpiochip0", gpiod.WithConsumer(fmt.Sprintf("aux%d", a.Number)))
	l, _ := c.RequestLine(cfg.GetInt(fmt.Sprintf("pin_aux%d", a.Number)), gpiod.AsOutput(1))
	time.Sleep(time.Second * 1)
	l.Close()
	c.Close()
}

func (a *Aux) Off() {
	cfg := config.Config()
	log.Infof("Aux %d off", a.Number)
	c, _ := gpiod.NewChip("gpiochip0", gpiod.WithConsumer(fmt.Sprintf("aux%d", a.Number)))
	l, _ := c.RequestLine(cfg.GetInt(fmt.Sprintf("pin_aux%d", a.Number)), gpiod.AsOutput(0))
	time.Sleep(time.Second * 1)
	l.Close()
	c.Close()
}

func (a *Aux) Trigger() {
	a.On()
	time.Sleep(time.Second * 1)
	a.Off()
}
