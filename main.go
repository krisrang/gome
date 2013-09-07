package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/krisrang/gome/drive"
)

const (
	VERSION      = "0.0.1"
	VERSIONFANCY = "Hairy Vermin"
)

var (
	port        = flag.String("port", "4000", "Port gome will run under")
	versionflag = flag.Bool("version", false, "Print version")
)

func main() {
	flag.Parse()

	if *versionflag {
		fmt.Println("Gome version", VERSION, VERSIONFANCY)
	} else {
		fmt.Println("Starting server on", *port)

		http.HandleFunc("/drive", drive.DriveHandler)

		err := http.ListenAndServe(":"+*port, nil)

		if err != nil {
			log.Fatal(err)
		}
	}
}
