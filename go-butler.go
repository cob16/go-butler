// A mumble bot based on the Gumble libary
// https://github.com/layeh/gumble/
package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/njdart/go-butler/configuration"
	"github.com/njdart/go-butler/steamgauge"
	"net"
	"regexp"
)

var (
	log          *logrus.Logger
	Steamconnect *regexp.Regexp
	ChatCommand  *regexp.Regexp
)

var comands = map[string]string{
	"help":   "Display this message.",
	"status": "show the status of steam",
}

//this is generated at init
var HelpString string

const HtmlNewLine string = `<br />`

//make to string to tell the user we haven't a clue
func ComandNotFound(usrInput string) string {
	return fmt.Sprintf("%s No command '%s' found! %s", HtmlNewLine, usrInput, HelpString)
}

func FormatHelpString(comands map[string]string) string {
	out := HtmlNewLine + "Available comands:"
	for cmd, helptext := range comands {
		out += fmt.Sprintf("%s <strong>%s</strong> - %s", HtmlNewLine, cmd, helptext)
	}
	return out
}

func init() {
	HelpString = FormatHelpString(comands)

	steamconnect, err := regexp.Compile(`connect ([A-Za-z0-9.:]+); *password (.*)`) //connect <ip>; <password>
	if err != nil {
		panic(err)
	}
	chatCmd, err := regexp.Compile(`!(\w+) *(.*)`) //separates 1st arg from the rest. first age must have a '!' e.g. !help arg1 arg2
	if err != nil {
		panic(err)
	}
	Steamconnect = steamconnect
	ChatCommand = chatCmd
}

// Create a html button to push that will allow joining
// into a multilayer source enging game server
func FormatSteamconnect(result []string) string {
	log.Info("steam link match ip: %d pass: %d", result[1], result[2])
	button := fmt.Sprintf("<br />IP: %s <br /> PASS: %s <br /><strong><a href='steam://connect/%s/%s'>CLICK TO CONNECT TO SERVER</a></strong><br />",
		result[1], result[2], result[1], result[2])
	log.Debug(button)
	return button
}

func HandleMessage(e *gumble.TextMessageEvent) {
	//parse steam connect strings and provide a html button to the channel
	result := Steamconnect.FindStringSubmatch(e.Message)
	if result != nil {
		e.Client.Self.Channel.Send(FormatSteamconnect(result), true)
	} else { //try user cmds instead
		result = ChatCommand.FindStringSubmatch(e.Message)
		if result != nil {
			switch result[1] {
			case "help":
				e.Sender.Send(HelpString)
			case "status":
				e.Client.Self.Channel.Send(SteamStatus(result), false)
			default:
				e.Sender.Send(ComandNotFound(result[1]))
			}
		}
	}
}

func SteamStatus(cmd []string) string {
	status, err := steamgauge.GetSteamStatus()
	if err != nil {
		panic(err)
	}
	if cmd[2] != "" {
		switch cmd[2] {
		case "tf2":
			return status.GetStatusTF2()
		case "csgo":
			return status.GetStatusCSGO()
		case "dota":
			return status.GetStatusDOTA2()
		default:
			return ComandNotFound(cmd[0])
		}
	}
	return status.GetStatus()
}

func main() {
	config, err := configuration.LoadConfiguration()
	if err != nil {
		panic(err)
	}
	log = config.GetLogger()
	log.Info("go-butler has sucessfully started!")
	tlsConfig, gumbleConfig := config.ExplodeConfiguration()

	keepAlive := make(chan bool)

	gumbleConfig.Attach(gumbleutil.Listener{
		UserChange: func(e *gumble.UserChangeEvent) {
			if e.Type.Has(gumble.UserChangeConnected) {
				e.User.Send("Welcome to the server, " + e.User.Name + "!")
			}
		},
		TextMessage: func(e *gumble.TextMessageEvent) {
			log.Infof("Received text message: %s\n", e.Message)
			HandleMessage(e)
		},
		//kill the program if we are disconnected
		Disconnect: func(e *gumble.DisconnectEvent) {
			log.Info("gumble was Disconnected from the server")
			keepAlive <- true
		},
	})

	log.Info("connecting to " + config.GetUri())
	client, err := gumble.DialWithDialer(new(net.Dialer), config.GetUri(), &gumbleConfig, &tlsConfig)
	if err != nil {
		log.Panic(err)
	} else {
		log.Info("connected!")
		log.Debug(client.State())
	}

	<-keepAlive
}
