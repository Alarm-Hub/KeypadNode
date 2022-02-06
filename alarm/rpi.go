// +build rpi

package alarm

import (
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/warthog618/gpiod"
  "time"
)

var c *gpiod.Chip

func (a *Alarm) Watch() {
	cfg := config.Config()
	c, _ = gpiod.NewChip("gpiochip0", gpiod.WithConsumer("AlarmWatcher"))
	l, err := c.RequestLine(cfg.GetInt("pin_alarm"), gpiod.AsInput)
	if err != nil {
    panic(err)
  }
	log.Debug("Watching for Alarm")
	for {
	  if i, _ := l.Value(); i == 0 {
	    time.Sleep(time.Second * 2)
	    if i, _ := l.Value(); i == 0 {
	      a.alarm()
      }
    }
    time.Sleep(time.Millisecond * 100)
  }
}

func (a *Alarm) alarm() {
	log.Info("Received alarm!")
	a.Emit("alarm", "true")
}
