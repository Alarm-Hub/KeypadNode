package alarm

import "github.com/eiannone/keyboard"

func (a *Alarm) Watch() {
	keyEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	for {
		event := <-keyEvents
		if event.Err != nil {
			panic(event.Err)
		}

		switch event.Key {
		case keyboard.KeySpace:
			a.Emit("alarm", "true")
		}
	}
}
