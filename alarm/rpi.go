// +build rpi

package alarm

import (
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/warthog618/gpiod"
)

var c *gpiod.Chip

func (a *Alarm) Watch() {
	cfg := config.Config()
	c, _ = gpiod.NewChip("gpiochip0", gpiod.WithConsumer("AlarmWatcher"))
	c.RequestLine(cfg.GetInt("pin_alarm"), gpiod.WithFallingEdge(a.alarm))
	log.Debug("Watching for Alarm")
}

func (a *Alarm) alarm(evt gpiod.LineEvent) {
	log.Info("Received alarm!")
	a.Emit("alarm", "true")
}
