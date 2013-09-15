package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/krisrang/gome/updater"
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

type PageData struct {
	Config *Config
}

type Config struct {
	Title       string
	GAID        string
	Description string
	Author      string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p := &PageData{Config: config}
	renderTemplate(w, "index.html", p)
}

func renderTemplate(w http.ResponseWriter, tpl string, data *PageData) {
	t, err := template.ParseFiles("templates/" + tpl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func loadConfig() {
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

	config = &conf

	fmt.Printf("Loaded config %v\n", conf)
}

func setupServer() {
	fmt.Println("Starting up http server on", *port)

	http.HandleFunc("/", mainPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.NotFoundHandler()

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

		loadConfig()
		go updater.SetupUpdater()
		setupServer()
	}
}
