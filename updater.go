package main

import (
	"fmt"
	"time"

	"github.com/krisrang/go-goodreads"
	"github.com/krisrang/go-lastfm"
	"github.com/krisrang/go-steam"

	"github.com/google/go-github/github"
)

var (
	LastTick time.Time

	LastfmUser   *lastfm.UserInfo
	LastfmTracks *[]lastfm.Track

	GithubUser  *github.User
	GithubRepos *[]github.Repository

	SteamUser  *steam.User
	SteamGames *steam.GamesList

	GoodreadsUser *goodreads.User
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
	goodreads.SetConfig(config.GoodreadsId, config.GoodreadsKey)
}

func runTimer() {
	tick := time.Tick(15 * time.Minute)
	for now := range tick {
		tock(now)
	}
}

func tock(now time.Time) {
	fmt.Println("Running update", now)

	GoodreadsUser = goodreads.GetUser()
	fmt.Println("Goodreads updated", time.Now())

	SteamUser = steam.GetUser()
	SteamGames = steam.GetGames()
	fmt.Println("Steam updated", time.Now())

	LastfmUser = lastfm.GetUser()
	LastfmTracks = lastfm.GetTracks(config.ClientLimit)
	fmt.Println("Last.fm updated", time.Now())

	GithubUser, GithubRepos = GithubUpdate(config.GithubToken, config.ClientLimit)
	fmt.Println("Github updated", time.Now())

	LastTick = time.Now()
	fmt.Println("Finished update", time.Now())
}
