package main

import (
	"flag"
	"fmt"
	"github.com/Phill93/DoorManager/log"
	"github.com/Phill93/DoorManager/version"
	"github.com/Phill93/DoorManager/wiegand"
	"time"
)

func handleEvents(events chan wiegand.Event) {
	for {
		e := <-events
		switch e.Type {
		case "key":
			log.Infof("Received key %s", e.Value)
		case "card":
			log.Infof("Received card id %s", e.Value)
		}
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "Version")
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

	defer wiegand.CleanGpios()
	go wiegand.InitGpio(events)
	go handleEvents(events)

	for {
		time.Sleep(1 * time.Second)
	}
}
