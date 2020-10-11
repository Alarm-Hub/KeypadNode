// +build rpi

package gate

import (
	"fmt"
	"github.com/Phill93/DoorManager/config"
	"github.com/warthog618/gpiod"
	"time"
)

func (g *Gate) open() {
	cfg := config.Config()
	fmt.Printf("Gate %s: ⇈\n\r", g.Name)
	c, _ := gpiod.NewChip("gpiochip0", gpiod.WithConsumer(g.Name))
	l, _ := c.RequestLine(cfg.GetInt("gate1_open"), gpiod.AsOutput(1))
	time.Sleep(time.Second * 1)
	l.SetValue(0)
	l.Close()
	c.Close()
}

func (g *Gate) close() {
	fmt.Printf("Gate %s: ⇊\n\n", g.Name)
}

func (g *Gate) stop() {
	fmt.Printf("Gate %s: ⇎\n\n", g.Name)
}
