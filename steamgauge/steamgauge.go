// go-lang steamgaug.es api wrapper and basic msg html generator
package steamgauge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Service struct {
	Online        int    `json:"online"` // 1 is up 2 is down
	Response_time int    `json:"time"`   // time in ms
	Error_msg     string `json:"error"`
}

//tf2 dota and CSG Service struct
type ValveGameService struct {
	TF2   Service `json:"440"` //the numbers are app id's
	DOTA2 Service `json:"570"`
	CSGO  Service `json:"730"`
}

//struct that holds all provided json data
type SteamStatus struct {
	Client struct {
		Online int `json:"online"`
	} `json:"ISteamClient"`
	Community   Service          `json:"SteamCommunity"`
	Store       Service          `json:"SteamStore"`
	User        Service          `json:"ISteamUser"`
	Items       ValveGameService `json:"IEconItems"`
	Matchmaking ValveGameService `json:"ISteamGameCoordinator"`
}

const apiv2url string = "https://steamgaug.es/api/v2"

const HtmlNewLine string = `<br />`

//makes the https request to steamgaug.es
// returns SteamStatus struct
func GetSteamStatus() (SteamStatus, error) {
	var status SteamStatus
	response, err := http.Get(apiv2url)
	if err != nil {
		return status, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return status, err
		}
		json.Unmarshal(contents, &status)
		return status, nil
	}
}

//get service.Online in bool
func (service Service) Bool() bool {
	if service.Online == 1 {
		return true
	} else {
		return false
	}
}

func HtmlColour(colour string, content string) string {
	return fmt.Sprintf("<span style='color:%s'>%s</span>", colour, content)
}

// format's service online status into a nice unicode green tick or red cross
// styled for mumble
func (service Service) FmtOnlineHtml() string {
	var tick, cross string = `☑`, `☒`
	var red, green = "#aa0000", "#00aa00"
	if service.Bool() {
		return HtmlColour(green, tick)
	} else {
		return HtmlColour(red, cross)
	}
}

//returns empty string if no error message
func filterForErrors(msg string) string {
	if msg != "No Error" {
		return msg
	}
	return ""
}

// formats html for status of provided Services
func getStatusGame(name string, items Service, matchmaking Service) string {
	return fmt.Sprintf("<strong>%s %s status:</strong> %s %s Item servers (%dms) %s  %s %s %s  %s",
		HtmlNewLine,
		name,
		HtmlNewLine,
		items.FmtOnlineHtml(),
		items.Response_time,
		filterForErrors(items.Error_msg),
		HtmlNewLine,
		matchmaking.FmtOnlineHtml(),
		" Matchmaking servers",
		filterForErrors(matchmaking.Error_msg),
	)
}

// formats html for status of provided Services
func (service Service) ServiceStatus(name string) string {
	return fmt.Sprintf("%s %s %s (%dms) %s",
		HtmlNewLine,
		service.FmtOnlineHtml(),
		name,
		service.Response_time,
		filterForErrors(service.Error_msg),
	)
}

//get the html formatting for the Client scruct
func (status SteamStatus) ClientOnlineHtml() string {
	checkbox := Service{Online: status.Client.Online}.FmtOnlineHtml()
	return fmt.Sprintln(checkbox, "Steam Client")
}

func (status SteamStatus) GetStatusSteam() string {
	return fmt.Sprintln(
		HtmlNewLine,
		"<strong>Steam Status:</strong>",
		HtmlNewLine,
		status.ClientOnlineHtml(),
		status.Community.ServiceStatus("Steam Community"),
		status.Store.ServiceStatus("Steam Store"),
		status.User.ServiceStatus("Steam User"),
	)
}

func (status SteamStatus) GetStatusTF2() string {
	return getStatusGame("Team Fortress 2", status.Items.TF2, status.Matchmaking.TF2)
}

func (status SteamStatus) GetStatusCSGO() string {
	return getStatusGame("Counter-Strike: Global Offensive", status.Items.CSGO, status.Matchmaking.CSGO)
}

func (status SteamStatus) GetStatusDOTA2() string {
	return getStatusGame("Defense of the Ancients", status.Items.DOTA2, status.Matchmaking.DOTA2)
}

//get the status of all services
func (status SteamStatus) GetStatus() string {
	return fmt.Sprintln(status.GetStatusSteam(), status.GetStatusTF2(), status.GetStatusCSGO(), status.GetStatusDOTA2())
}
