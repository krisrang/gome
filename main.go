package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

const (
	VERSION      = "0.0.1"
	VERSIONFANCY = "Hairy Vermin"
)

var (
	port        = flag.String("port", "4000", "Port gome will run under")
	versionflag = flag.Bool("version", false, "Print version")
)

type PageData struct {
	Title string
}

func setupUpdater() {
	tick := time.Tick(15 * time.Minute)
	for now := range tick {
		fmt.Printf("%v\n", now)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	p := &PageData{Title: "test"}
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

func main() {
	flag.Parse()

	if *versionflag {
		fmt.Println("Gome version", VERSION, VERSIONFANCY)
	} else {
		fmt.Println("Setting up data updater and running first run")
		go setupUpdater()

		fmt.Println("Starting server on", *port)

		http.HandleFunc("/", mainPage)
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
		http.NotFoundHandler()

		err := http.ListenAndServe(":"+*port, nil)

		if err != nil {
			log.Fatal(err)
		}
	}
}
