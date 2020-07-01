package main

import (
	"flag"
	"fmt"
	"github.com/Phill93/DoorManager/version"
	"github.com/Phill93/DoorManager/wiegand"
	"time"
)

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
	fmt.Println("Hello.")
	wiegand.InitGpio()
	time.Sleep(60 * time.Second)
	defer wiegand.CleanGpios()
}
