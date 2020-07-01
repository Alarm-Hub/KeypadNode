package wiegand

import (
	"fmt"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var c *gpiod.Chip
var low uint32 = 0x0
var high uint32 = 0x0
var carry int = 0

//Data
func lowHandler(evt gpiod.LineEvent) {
	fmt.Println("Low!")
	low = low ^ (1 << carry)
	carry += 1
	fmt.Printf("%26b\n", low)
}

//Parity
func highHandler(evt gpiod.LineEvent) {
	fmt.Println("High!")
	high = high ^ (1 << carry)
	carry += 1
	fmt.Printf("%26b\n", high)
}

func InitGpio() {
	c, _ = gpiod.NewChip("gpiochip0", gpiod.WithConsumer("DoorManager_Wiegand"))
	c.RequestLine(rpi.GPIO14, gpiod.WithFallingEdge(lowHandler))
	c.RequestLine(rpi.GPIO15, gpiod.WithFallingEdge(highHandler))
}

func CleanGpios() {
	c.Close()
}
