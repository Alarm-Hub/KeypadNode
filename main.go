package main

import (
	"flag"
	"fmt"
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	"github.com/Phill93/DoorManager/version"
	"github.com/Phill93/DoorManager/wiegand"
	"time"
)

func checkTimeout(timestamp time.Time, timeout int) bool {
	now := time.Now()
	offset := now.Sub(timestamp)
	if offset.Seconds() > float64(timeout) {
		log.Infof("timeout reached! %d", offset.Seconds())
		return true
	} else {
		log.Infof("timeout not reached! %d", offset.Seconds())
		return false
	}
}

func handleEvents(events chan wiegand.Event) {
	var code []string
	var timestamp time.Time
	for {
		e := <-events
		switch e.Type {
		case "key":
			log.Infof("Received key %s", e.Value)
			if e.Value != "ENT" && e.Value != "ESC" {
				if len(code) == 0 {
					timestamp = time.Now()
					log.Info("Recived key %v", timestamp)
				}
				code = append(code, e.Value)
			} else if e.Value == "ENT" {
				log.Infof("code is %s", code)
				code = nil
			} else if e.Value == "ESC" || checkTimeout(timestamp, 2) {
				log.Info("clear code buffer")
				code = nil
			}
		case "card":
			log.Infof("Received card id %s", e.Value)
		}
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "Version")
	c := config.Config()
	flag.Parse()
	if *versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		return
	}

	events := make(chan wiegand.Event, 1)

	fmt.Print(c.Get("Test"))
	defer wiegand.CleanGpios()
	go wiegand.InitGpio(events)
	go handleEvents(events)

	for {
		time.Sleep(1 * time.Second)
	}
}
