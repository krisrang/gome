package updater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	LastTick time.Time
)

func SetupUpdater(c *main.Config) {
	fmt.Println("Setting up data updater and running first run")

	tock(time.Now())

	tick := time.Tick(15 * time.Minute)
	for now := range tick {
		tock(now)
	}
}

func tock(now time.Time) {
	LastTick = now
	fmt.Println("Running update", now)

	LastfmUpdate(now)
}

func getRequest(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func jsonUnmarshal(b []byte, i interface{}) {
	err := json.Unmarshal(b, i)
	if err != nil {
		log.Fatal(err)
	}
}
