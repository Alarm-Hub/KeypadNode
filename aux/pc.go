// +build !rpi

package aux

import (
	"github.com/Phill93/DoorManager/log"
	"time"
)

func (a *Aux) On() {
	log.Infof("Aux %s on", a.Number)
}

func (a *Aux) Off() {
	log.Infof("Aux %s off", a.Number)
}

func (a *Aux) Trigger() {
	a.on()
	time.Sleep(time.Second * 1)
	a.off()
}
