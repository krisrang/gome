package updater

import (
	"fmt"
	"time"
	// "log"
)

var (
	LastfmUserData LastfmUser // primary export

	apiRoot = "http://ws.audioscrobbler.com/2.0"
)

type LastfmUser struct {
	User LastfmUserInfo
}

type LastfmUserInfo struct {
	Name      string
	Realname  string
	URL       string
	PlayCount string
	Image     []LastfmImage
}

type LastfmImage struct {
	Text string `json:"#text"`
	Size string
}

func LastfmUpdate(now time.Time) {
	uri := apiRoot + "?method=user.getinfo&user=noin&format=json&api_key=21b8bc82727c669ccaa55bbe864b1fff"
	LastfmUserData = LastfmUser{}

	data := getRequest(uri)
	jsonUnmarshal(data, &LastfmUserData)

	fmt.Println("Last.fm updated", now)
}
