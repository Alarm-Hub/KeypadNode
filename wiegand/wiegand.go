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
var low uint32 = 0x0
var high uint32 = 0x0
var carry uint32 = 0
var lowWd chan bool
var highWd chan bool
var end = false

type Keypad struct {
	listeners map[string][]chan string
}

func (k *Keypad) AddListener(e string, ch chan string) {
	if k.listeners == nil {
		k.listeners = make(map[string][]chan string)
	}
	log.Debugf("Adding channel for event %s", e)
	if _, ok := k.listeners[e]; ok {
		k.listeners[e] = append(k.listeners[e], ch)
	} else {
		k.listeners[e] = []chan string{ch}
	}
}

func (k *Keypad) RemoveListener(e string, ch chan string) {
	log.Debugf("Removing channel for event %s", e)
	if _, ok := k.listeners[e]; ok {
		for i := range k.listeners[e] {
			if k.listeners[e][i] == ch {
				k.listeners[e] = append(k.listeners[e][:i], k.listeners[e][i+1:]...)
				break
			}
		}
	}
}

func (k *Keypad) Emit(e string, response string) {
	log.Debugf("Event %s is emitted!", e)
	if _, ok := k.listeners[e]; ok {
		for _, handler := range k.listeners[e] {
			fmt.Print("Found!")
			go func(handler chan string) {
				log.Debug("Writing data!")
				handler <- response
			}(handler)
		}
	}
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

func InitReader(pad *Keypad) {
	cfg := config.Config()
	c, _ = gpiod.NewChip("gpiochip0", gpiod.WithConsumer("KeypadNode"))
	lowWd = make(chan bool, 1)
	c.RequestLine(cfg.GetInt("data_low"), gpiod.WithFallingEdge(lowHandler))
	highWd = make(chan bool, 1)
	c.RequestLine(cfg.GetInt("data_high"), gpiod.WithFallingEdge(highHandler))
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
