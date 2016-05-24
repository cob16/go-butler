// golang steamgaug.es api wrapper
package steamgauge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
		online int `json:"online"`
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

//format's service online status into a nice unicode green tick or red cross
func (service Service) FmtOnlineHtml() string {
	var tick, cross string = `☑`, `☒`
	var red, green = "#aa0000", "#00aa00"
	if service.Bool() {
		return HtmlColour(green, tick)
	} else {
		return HtmlColour(red, cross)
	}
}

// formats html for status of provided Services
func getStatus(name string, items Service, matchmaking Service) string {
	return fmt.Sprintf("<strong>%s %s status:</strong> %s %s Item servers (%dms) %s  %s %s %s  %s",
		HtmlNewLine,
		name,
		HtmlNewLine,
		items.FmtOnlineHtml(),
		items.Response_time,
		items.Error_msg,
		HtmlNewLine,
		matchmaking.FmtOnlineHtml(),
		" Matchmaking servers",
		matchmaking.Error_msg,
	)
}

func (status SteamStatus) GetStatusTF2() string {
	return getStatus("Team Fortress 2", status.Items.TF2, status.Matchmaking.TF2)
}

func (status SteamStatus) GetStatusCSGO() string {
	return getStatus("Counter-Strike: Global Offensive", status.Items.CSGO, status.Matchmaking.CSGO)
}

func (status SteamStatus) GetStatusDOTA2() string {
	return getStatus("Defense of the Ancients", status.Items.DOTA2, status.Matchmaking.DOTA2)
}

func (status SteamStatus) GetStatus() string {
	return strings.Join([]string{status.GetStatusTF2(), status.GetStatusCSGO(), status.GetStatusDOTA2()}, "")
}
