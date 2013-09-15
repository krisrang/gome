package updater

import (
	"fmt"
	"time"
)

func SetupUpdater() {
	fmt.Println("Setting up data updater and running first run")

	tick := time.Tick(15 * time.Minute)
	for now := range tick {
		fmt.Printf("%v\n", now)
	}
}
