package wiegand

import (
	"fmt"
	"github.com/Phill93/DoorManager/log"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"math"
	"math/bits"
	"time"
)

var c *gpiod.Chip
var low uint32 = 0x0
var high uint32 = 0x0
var carry uint32 = 0
var lowWd chan bool
var highWd chan bool
var end = false
var channel chan Event

type Event struct {
	Type  string
	Value string
}

//Data
func lowHandler(evt gpiod.LineEvent) {
	low = low ^ (1 << carry)
	carry += 1
	lowWd <- true
}

//Parity
func highHandler(evt gpiod.LineEvent) {
	high = high ^ (1 << carry)
	carry += 1
	highWd <- true
}

func reverse(x uint32, size uint32) uint32 {
	return bits.Reverse32(x) >> (32 - size)
}

func parseData(data uint32) string {
	if data == 10 {
		return "ESC"
	} else if data == 11 {
		return "ENT"
	} else {
		return fmt.Sprint(data)
	}
}

func processData(data uint32, parity uint32, bits uint32) {
	data = reverse(data, bits)
	parity = reverse(parity, bits)
	if (data + parity) == uint32(math.Pow(2, float64(bits)))-1 {
		log.Info("Parity ok!")
		pdata := parseData(data)
		if bits == 4 {
			channel <- Event{
				Type:  "key",
				Value: pdata,
			}
		} else if bits == 26 {
			channel <- Event{
				Type:  "card",
				Value: pdata,
			}
		}
	}
}

func InitGpio(ch chan Event) {
	channel = ch
	c, _ = gpiod.NewChip("gpiochip0", gpiod.WithConsumer("DoorManager_Wiegand"))
	lowWd = make(chan bool, 1)
	c.RequestLine(rpi.GPIO14, gpiod.WithFallingEdge(lowHandler))
	highWd = make(chan bool, 1)
	c.RequestLine(rpi.GPIO15, gpiod.WithFallingEdge(highHandler))
	for {
		select {
		case <-lowWd:
			log.Infof("Data received on low! (c: %d)\n", carry)
		case <-highWd:
			log.Infof("Data received on high! (c: %d)\n", carry)
		case <-time.After(4000 * time.Microsecond):
			if carry == 4 || carry == 26 {
				processData(low, high, carry)
				low = 0
				high = 0
				carry = 0
			} else if carry > 0 {
				low = 0
				high = 0
				carry = 0
			}
		}
		if end {
			break
		}
	}
}

func CleanGpios() {
	end = true
	c.Close()
}
