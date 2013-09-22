package clients

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	SteamUserData  *SteamUser
	SteamGamesData *SteamGamesList

	steamApiRoot = "http://steamcommunity.com/id/"
)

type SteamUser struct {
	SteamID64      string `xml:"steamID64"`
	SteamID        string `xml:"steamID"`
	OnlineState    string `xml:"onlineState"`
	StateMessage   string `xml:"stateMessage"`
	AvatarIcon     string `xml:"avatarIcon"`
	AvatarMedium   string `xml:"avatarMedium"`
	AvatarFull     string `xml:"avatarFull"`
	CustomURL      string `xml:"customURL"`
	MemberSince    string `xml:"memberSince"`
	SteamRating    string `xml:"steamRating"`
	HoursPlayed2Wk string `xml:"hoursPlayed2Wk"`
	Location       string `xml:"location"`
	Realname       string `xml:"realname"`
	Summary        string `xml:"summary"`
	GameCount      int
	// MostPlayed     []SteamMostPlayedGame `xml:"mostPlayedGames>mostPlayedGame"`
}

func (u SteamUser) FullURL() string {
	return steamApiRoot + u.CustomURL
}

func (u SteamUser) RatingDescription() string {
	title := "Playing on PS3"

	switch u.SteamRating {
	case "10":
		title = "EAGLES SCREAM"
	case "9":
		title = "Still not 10"
	case "8":
		title = "COBRA KAI!"
	case "7":
		title = "Wax on, Wax off"
	case "6":
		title = "Oooh! Shiny!"
	case "5":
		title = "Halfway Cool"
	case "4":
		title = "Master of Nothing"
	case "3":
		title = "Shooting Blanks"
	case "2":
		title = "Nearly Lifeless"
	case "1":
		title = "El Terrible!"
	}

	return title
}

// type SteamMostPlayedGame struct {
// 	Name          string `xml:"gameName"`
// 	Link          string `xml:"gameLink"`
// 	Icon          string `xml:"gameIcon"`
// 	Logo          string `xml:"gameLogo"`
// 	LogoSmall     string `xml:"gameLogoSmall"`
// 	HoursPlayed   string `xml:"hoursPlayed"`
// 	HoursOnRecord string `xml:"hoursOnRecord"`
// }

type SteamGamesList struct {
	Games []SteamGame `xml:"games>game"`
}

type SteamGame struct {
	AppID           string `xml:"appID"`
	Name            string `xml:"name"`
	Logo            string `xml:"logo"`
	StoreLink       string `xml:"storeLink"`
	HoursLast2Weeks string `xml:"hoursLast2Weeks"`
	HoursOnRecord   string `xml:"hoursOnRecord"`
}

type SteamGamesByLast2Weeks SteamGamesList
type SteamGamesByHours SteamGamesList

// yes, this is bad
func (s SteamGamesByLast2Weeks) Less(i, j int) bool {
	a, _ := strconv.ParseFloat(s.Games[i].HoursLast2Weeks, 64)
	b, _ := strconv.ParseFloat(s.Games[j].HoursLast2Weeks, 64)

	return a > b
}

func (s SteamGamesByHours) Less(i, j int) bool {
	a, _ := strconv.ParseFloat(s.Games[i].HoursOnRecord, 64)
	b, _ := strconv.ParseFloat(s.Games[j].HoursOnRecord, 64)

	return a > b
}

func (s SteamGamesByLast2Weeks) Len() int      { return len(s.Games) }
func (s SteamGamesByLast2Weeks) Swap(i, j int) { s.Games[i], s.Games[j] = s.Games[j], s.Games[i] }

func (s SteamGamesByHours) Len() int      { return len(s.Games) }
func (s SteamGamesByHours) Swap(i, j int) { s.Games[i], s.Games[j] = s.Games[j], s.Games[i] }

func SteamUpdate(user string) {
	steamUserUpdate(user)
	steamGamesUpdate(user)

	SteamUserData.Summary = strings.Replace(SteamUserData.Summary, "<br>", "", -1)
	SteamUserData.GameCount = len(SteamGamesData.Games)

	playedLast2Weeks := []SteamGame{}
	notPlayedLastWeeks := []SteamGame{}

	for i := range SteamGamesData.Games {
		game := SteamGamesData.Games[i]

		if game.HoursLast2Weeks != "" {
			playedLast2Weeks = append(playedLast2Weeks, game)
		} else {
			notPlayedLastWeeks = append(notPlayedLastWeeks, game)
		}
	}

	sort.Sort(SteamGamesByLast2Weeks{Games: playedLast2Weeks})
	sort.Sort(SteamGamesByHours{Games: notPlayedLastWeeks})

	sortedGames := append(playedLast2Weeks, notPlayedLastWeeks...)

	SteamGamesData.Games = sortedGames[:5]

	fmt.Println("Steam updated", time.Now())
}

func steamUserUpdate(user string) {
	uri := steamApiRoot + user + "?xml=1"
	SteamUserData = &SteamUser{}

	data := getRequest(uri)
	xmlUnmarshal(data, SteamUserData)
}

func steamGamesUpdate(user string) {
	uri := steamApiRoot + user + "/games/?xml=1"
	SteamGamesData = &SteamGamesList{}

	data := getRequest(uri)
	xmlUnmarshal(data, SteamGamesData)
}
