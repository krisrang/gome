package main

import (
	"fmt"
	"time"

	"github.com/krisrang/gome/clients"
)

var (
	LastTick time.Time
)

func setupUpdater() {
	fmt.Println("Setting up data updater and running first run")
	tock(time.Now())
	go setupTimer()
}

func setupTimer() {
	tick := time.Tick(15 * time.Minute)
	for now := range tick {
		tock(now)
	}
}

func tock(now time.Time) {
	fmt.Println("Running update", now)

	clients.SteamUpdate(config.SteamUser)
	clients.LastfmUpdate(config.LastfmUser, config.LastfmKey)
	clients.GithubUpdate(config.GithubToken)

	LastTick = time.Now()
	fmt.Println("Finished update", time.Now())
}
