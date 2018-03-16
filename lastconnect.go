package main

import (
	"fmt"
	"layeh.com/gumble/gumble"
	"time"
)

var connect = &Command{
	Run:            Lastconnect,
	PublicResponse: false,
	UsageLine:      "connect [delete]",
	Short:          "Shows last connect info pasted into mumble",
	Long: `
Shows last connect info last pasted into mumbble, you can allso
use the "delete" arg to remove this from memory.

note: This will replace the connect string with your user ID!
`,
}

type SourceConnect struct {
	hostname string
	password string
	UserID   uint32
}

// Create a html button to push that will allow joining
// into a multilayer source enging game server
func (c SourceConnect) HTMLSteamConnect() string {
	timeStamp := time.Now().UTC().Format(time.RFC850)
	button := fmt.Sprintf("<br />Here is the last connect posted:<br />[%s]<pre>%s</pre><strong><a title='This will connect to the server via steam' href='steam://connect/%s/%s'>CLICK TO CONNECT TO SERVER</a></strong>",
		timeStamp, c.ConnectString(), c.hostname, c.password)
	log.Debug(button)
	return button
}

func (c SourceConnect) ConnectString() string {
	return fmt.Sprintf("connect %s; password %s", c.hostname, c.password)
}

var currentConnect = new(SourceConnect)
var currentConnectHTML = "<strong>There is not yet a connect!</strong>"
var currentConnectString = "There is not yet a connect!"

func (c SourceConnect) GenConnectString() (string, string) {
	return c.HTMLSteamConnect(), c.ConnectString()
}

func Lastconnect(cmd *Command, args []string, event *gumble.TextMessageEvent) string {
	sender := event.Sender
	if args[2] != "" {
		switch args[2] {
		case "delete":
			var delMsg = fmt.Sprintf("The last connect was deleted by '%s' ID: %d at %s",
				sender.Name, sender.UserID, time.Now().UTC().Format(time.RFC850))

			currentConnectHTML = delMsg
			currentConnectString = delMsg
			currentConnect = new(SourceConnect)
		case "raw":
			return currentConnectString
		default:
			return CommandNotFound(args[0])
		}
	}
	return currentConnectHTML
}
