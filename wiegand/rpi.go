// +build rpi

package wiegand

import (
	"fmt"
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/warthog618/gpiod"
	"math"
	"math/bits"
	"time"
)

var c *gpiod.Chip

func InitReader(pad *Keypad) {
	cfg := config.Config()
	c, _ = gpiod.NewChip("gpiochip0", gpiod.WithConsumer("KeypadNode"))
	lowWd = make(chan bool, 1)
	c.RequestLine(cfg.GetInt("data_low"), gpiod.WithFallingEdge(lowHandler))
	highWd = make(chan bool, 1)
	c.RequestLine(cfg.GetInt("data_high"), gpiod.WithFallingEdge(highHandler))
	defer CleanGpios()
	time.Sleep(time.Second)
	for {
		select {
		case <-lowWd:
			log.Debugf("Data received on low! (c: %d)\n", carry)
		case <-highWd:
			log.Debugf("Data received on high! (c: %d)\n", carry)
		case <-time.After(4000 * time.Microsecond):
			if carry == 4 || carry == 26 {
				pad.processData(low, high, carry)
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

func (k *Keypad) processData(data uint32, parity uint32, bits uint32) {
	data = reverse(data, bits)
	parity = reverse(parity, bits)
	if (data + parity) == uint32(math.Pow(2, float64(bits)))-1 {
		log.Debug("Parity ok!")
		pdata := parseData(data)
		if bits == 4 {
			k.Emit("key", pdata)
		} else if bits == 26 {
			k.Emit("card", pdata)
		}
	}
}
