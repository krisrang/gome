package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/krisrang/go-goodreads"
	"github.com/krisrang/go-lastfm"
	"github.com/krisrang/go-steam"

	"github.com/google/go-github/github"
)

const (
	VERSION      = "0.1.0"
	VERSIONFANCY = "Jovial Gopher"
)

var (
	versionflag = flag.Bool("version", false, "Print version")

	config      *Config
	currentData = &PageData{AllSynced: false}
)

type Config struct {
	GAID          string
	ClientLimit   int
	UpdateMinutes int

	LastfmUser string
	LastfmKey  string

	GithubToken string

	SteamUser string

	GoodreadsId  string
	GoodreadsKey string
}

type PageData struct {
	Config    *Config
	AllSynced bool

	LastfmUser   *lastfm.UserInfo
	LastfmTracks *[]lastfm.Track

	GithubUser  *github.User
	GithubRepos *[]github.Repository

	SteamUser  *steam.User
	SteamGames *steam.GamesList

	GoodreadsUser *goodreads.User
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	renderTemplate(w, "index.html", currentData)
}

func statusPage(w http.ResponseWriter, r *http.Request) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	fmt.Fprintln(w, "RAM: used", m.Alloc/1024, "allocated", m.Sys/1024)
	fmt.Fprintln(w, "Last updater tick:", LastTick)
}

func renderTemplate(w http.ResponseWriter, tpl string, data interface{}) {
	t, err := template.ParseFiles("templates/"+tpl,
		"templates/lastfm.html", "templates/github.html",
		"templates/steam.html", "templates/goodreads.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func loadConfig() *Config {
	fmt.Println("Loading configuration")

	conf := Config{}

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err)
	}

	conf.GAID = os.Getenv("GAID")
	conf.LastfmUser = os.Getenv("LASTFMUSER")
	conf.LastfmKey = os.Getenv("LASTFMKEY")
	conf.GithubToken = os.Getenv("GITHUBTOKEN")
	conf.SteamUser = os.Getenv("STEAMUSER")
	conf.GoodreadsId = os.Getenv("GOODREADSID")
	conf.GoodreadsKey = os.Getenv("GOODREADSKEY")

	return &conf
}

func setupServer() {
	port := os.Getenv("PORT")
	fmt.Println("Starting up http server on", port)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/status", statusPage)
	http.HandleFunc("/", mainPage)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *versionflag {
		fmt.Println("Gome version", VERSION, VERSIONFANCY)
	} else {
		fmt.Println("Gome version", VERSION)

		config = loadConfig()
		go setupUpdater()
		setupServer()
	}
}
