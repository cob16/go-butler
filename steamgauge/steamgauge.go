// golang steamgaug.es api wrapper
package steamgauge

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type Service struct {
	Online        int `json:"online"` // 1 is up 2 is down
	Response_time int `json:"time"` // time in ms
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
	Community Service `json:"SteamCommunity"`
	Store Service `json:"SteamStore"`
	User Service `json:"ISteamUser"`
  Items ValveGameService `json:"IEconItems"`
	Matchmaking ValveGameService`json:"ISteamGameCoordinator"`
}

const apiv2url string = "https://steamgaug.es/api/v2"

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

//format's service into a nice html string
func (service Service) FmtOnlineHtml() (string){
	var tick, cross string = `☑`, `☒`
	if service.Online == 1 {
		return tick
	} else {
		return cross
	}
}

