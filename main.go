package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"html/template"
	"runtime"

	"github.com/krisrang/gome/clients"
)

const (
	VERSION      = "0.0.1"
	VERSIONFANCY = "Hairy Vermin"
)

var (
	port        = flag.String("port", "4000", "Port gome will run under")
	versionflag = flag.Bool("version", false, "Print version")

	config *Config
)

type Config struct {
	GAID string

	LastfmUser string
	LastfmKey  string
}

type PageData struct {
	Config       *Config
	LastfmUser   *clients.LastfmUserInfo
	LastfmTracks *[]clients.LastfmTrack
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p := &PageData{
		Config:       config,
		LastfmUser:   &clients.LastfmUserData.User,
		LastfmTracks: &clients.LastfmTrackData.Tracks.Tracks,
	}
	renderTemplate(w, "index.html", p)
}

func statusPage(w http.ResponseWriter, r *http.Request) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	// fmt.Fprintln(w, "PID:", os.Getpid())
	fmt.Fprintln(w, "RAM: used", m.Alloc/1024, "allocated", m.Sys/1024)
	fmt.Fprintln(w, "Last updater tick:", LastTick)
}

func renderTemplate(w http.ResponseWriter, tpl string, data *PageData) {
	t, err := template.ParseFiles("templates/" + tpl)

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

	return &conf

	// fmt.Printf("Loaded config %v\n", conf)
}

func setupServer() {
	fmt.Println("Starting up http server on", *port)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/status", statusPage)
	http.HandleFunc("/", mainPage)

	err := http.ListenAndServe(":"+*port, nil)

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
