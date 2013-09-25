package main

import (
	"fmt"
	"time"

	"github.com/krisrang/go-goodreads"
	"github.com/krisrang/go-lastfm"
	"github.com/krisrang/go-steam"
)

var (
	LastTick time.Time
)

func setupUpdater() {
	fmt.Println("Setting up data updater and running first run")

	setup()
	tock(time.Now())

	go runTimer()
}

func setup() {
	lastfm.SetConfig(config.LastfmUser, config.LastfmKey)
	steam.SetConfig(config.SteamUser, config.ClientLimit)
}

func runTimer() {
	fmt.Println("Updating every", config.UpdateMinutes, "minutes")
	tick := time.Tick(time.Duration(config.UpdateMinutes) * time.Minute)
	for now := range tick {
		tock(now)
	}
}

func tock(now time.Time) {
	fmt.Println("Running update", now)

	goodreadsuser := goodreads.GetUser(config.GoodreadsId, config.GoodreadsKey, config.ClientLimit)
	fmt.Println("Goodreads updated", time.Now())

	steamuser := steam.GetUser()
	steamgames := steam.GetGames()
	fmt.Println("Steam updated", time.Now())

	lastfmuser := lastfm.GetUser()
	lastfmtracks := lastfm.GetTracks(config.ClientLimit)
	fmt.Println("Last.fm updated", time.Now())

	githubuser, githubrepos := GithubUpdate(config.GithubToken, config.ClientLimit)
	fmt.Println("Github updated", time.Now())

	currentData = &PageData{
		Config:        config,
		LastfmUser:    lastfmuser,
		LastfmTracks:  lastfmtracks,
		GithubUser:    githubuser,
		GithubRepos:   githubrepos,
		SteamUser:     steamuser,
		SteamGames:    steamgames,
		GoodreadsUser: goodreadsuser,
		AllSynced:     true,
	}

	LastTick = time.Now()

	fmt.Println("Finished update", time.Now())
}
