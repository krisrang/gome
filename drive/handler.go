package drive

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/drive/v2"
)

var config = &oauth.Config{
	ClientId:     os.Getenv("GAPI_KEY"),
	ClientSecret: os.Getenv("GAPI_SECRET"),
	Scope:        "https://www.googleapis.com/auth/drive",
	RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

var cacheFile = "tokenCache"

func DriveHandler(w http.ResponseWriter, r *http.Request) {
	service, _ := drive.New(getOAuthClient(config))

	files, _ := allFiles(service)

	for i := range files {
		file := files[i]
		fmt.Fprintln(w, file.OriginalFilename)
	}
}

func allFiles(d *drive.Service) ([]*drive.File, error) {
	var fs []*drive.File
	pageToken := ""
	for {
		q := d.Files.List()
		// If we have a pageToken set, apply it to the query
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			return fs, err
		}
		fs = append(fs, r.Items...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return fs, nil
}

func RequestToken() {
	authUrl := config.AuthCodeURL("state")
	fmt.Printf("Go to the following link in your browser: %v\n", authUrl)
	t := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}

	// Read the code, and exchange it for a token.
	fmt.Printf("Enter verification code: ")
	var code string
	fmt.Scanln(&code)

	if _, err := t.Exchange(code); err != nil {
		log.Fatal(err)
	}

	saveToken(cacheFile, t.Token)
}

func getOAuthClient(config *oauth.Config) *http.Client {
	token, err := tokenFromFile(cacheFile)

	if err != nil {
		log.Fatal(err)
	}

	t := &oauth.Transport{
		Token:     token,
		Config:    config,
		Transport: http.DefaultTransport,
	}
	return t.Client()
}

func tokenFromFile(file string) (*oauth.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(file string, token *oauth.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}
