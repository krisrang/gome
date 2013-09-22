package clients

import (
	"fmt"
	"time"
)

var (
	LastfmUserData  *LastfmUser   // primary export
	LastfmTrackData *LastfmTracks // primary export

	lastfmApiRoot = "http://ws.audioscrobbler.com/2.0"
)

type LastfmUser struct {
	User LastfmUserInfo
}

type LastfmUserInfo struct {
	Name       string
	Realname   string
	URL        string
	PlayCount  string
	Country    string
	Image      []LastfmImage
	Registered LastfmDate
}

func (t LastfmUserInfo) GetImage() string {
	image := t.Image[len(t.Image)-1]
	return image.URL
}

type LastfmTracks struct {
	Tracks LastfmTackList `json:"recenttracks"`
}

type LastfmTackList struct {
	Tracks []LastfmTrack `json:"track"`
}

type LastfmTrack struct {
	Artist LastfmArtist
	Name   string
	URL    string
	MBID   string
	NP     LastfmNowPlaying `json:"@attr"`
	Image  []LastfmImage
	Date   LastfmDate
}

func (t LastfmTrack) IsNowPlaying() bool {
	return t.NP.NowPlaying == "true"
}

type LastfmNowPlaying struct {
	NowPlaying string
}

func (t LastfmTrack) GetImage() string {
	image := t.Image[len(t.Image)-1]
	return image.URL
}

type LastfmArtist struct {
	Name string `json:"#text"`
	MBID string
}

type LastfmImage struct {
	URL  string `json:"#text"`
	Size string
}

type LastfmDate struct {
	Text string `json:"#text"`
	UTS  string `json:"uts,unixtime"`
}

func (d LastfmDate) ParseDate() (time.Time, error) {
	date, err := time.Parse("2006-01-02 15:04", d.Text)
	if err != nil {
		date, err = time.Parse("02 Jan 2006, 15:04", d.Text)
		if err != nil {
			return time.Time{}, err
		}
	}

	return date, nil
}

func (d LastfmDate) ShortDate() string {
	date, err := d.ParseDate()
	if err != nil {
		return ""
	}

	return (string)(date.Format("2 Jan 2006"))
}

func (d LastfmDate) RelativeDate() string {
	date, err := d.ParseDate()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	s := time.Now().Sub(date)

	days := int(s / (24 * time.Hour))
	if days > 1 {
		return fmt.Sprintf("%v days ago", days)
	} else if days == 1 {
		return fmt.Sprintf("%v day ago", days)
	}

	hours := int(s / time.Hour)
	if hours > 1 {
		return fmt.Sprintf("%v hours ago", hours)
	}

	minutes := int(s / time.Minute)
	if minutes > 2 {
		return fmt.Sprintf("%v minutes ago", minutes)
	} else {
		return "Just now"
	}
}

func LastfmUpdate(user string, key string) {
	lastfmUserUpdate(user, key)
	lastfmTracksUpdate(user, key)
	fmt.Println("Last.fm updated", time.Now())
}

func lastfmUserUpdate(user string, key string) {
	uri := lastfmApiRoot + "?method=user.getinfo&format=json&user=" + user + "&api_key=" + key
	LastfmUserData = &LastfmUser{}

	data := getRequest(uri)
	jsonUnmarshal(data, LastfmUserData)
}

func lastfmTracksUpdate(user string, key string) {
	uri := lastfmApiRoot + "?method=user.getrecenttracks&format=json&user=" + user + "&api_key=" + key
	LastfmTrackData = &LastfmTracks{}

	data := getRequest(uri)
	jsonUnmarshal(data, LastfmTrackData)

	LastfmTrackData.Tracks.Tracks = LastfmTrackData.Tracks.Tracks[:5]
}
