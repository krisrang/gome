package clients

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

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
