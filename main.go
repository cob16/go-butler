// A mumble bot based on the Gumble libary
// https://github.com/layeh/gumble/
package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/njdart/go-butler/configuration"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	log          *logrus.Logger
	Steamconnect *regexp.Regexp
	ChatCommand  *regexp.Regexp
)

//taken from cmd/go/main as it can not be imported
// A Command is an implementation of command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) string

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the '!help' output.
	Short string

	// Long is the long message shown in the !help <command> output.
	Long string

	//if the cmd output will be sent back to the whole channel
	// or to a user in a private message
	PublicResponse bool
}

//yes more duplication as that package is not importable
// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

var commands = []*Command{
	status,
}

//this is generated at init
var HelpString string

var help = &Command{
	PublicResponse: false,
	UsageLine:      "help [command]",
	Short:          "cmds list as well as cmd deail",
	Long:           ``,
}

func FormatHelpString(commands []*Command) string {
	commands = append(commands, help)
	out := HtmlNewLine + "Available comands:"
	for _, cmd := range commands {
		out += fmt.Sprintf("%s <strong>%s</strong> - %s", HtmlNewLine, cmd.Name(), cmd.Short)
	}
	return out
}

func FormatCmdHelp(cmd *Command) string {
	return fmt.Sprintln(cmd.UsageLine, HtmlNewLine, cmd.Long)
}

//called with special case case in HandleMessage
func Help(args []string) string {
	if args[2] == "" { //just help
		return HelpString
	}
	//look for other cmds passed as args
	// and show there UsageLine + Long discription
	for _, cmd := range commands {
		if cmd.Name() == args[2] {
			return FormatCmdHelp(cmd)
		}
	}
	return CommandNotFound(args[2])
}

const HtmlNewLine string = `<br />`

//make to string to tell the user we haven't a clue
func CommandNotFound(usrInput string) string {
	return fmt.Sprintf("%s No command '%s' found! %s", HtmlNewLine, usrInput, HelpString)
}

// Create a html button to push that will allow joining
// into a multilayer source enging game server
func FormatSteamConnect(result []string) string {
	log.Info("steam link match ip: %d pass: %d", result[1], result[2])
	button := fmt.Sprintf("<br />IP: %s <br /> PASS: %s <br /><strong><a href='steam://connect/%s/%s'>CLICK TO CONNECT TO SERVER</a></strong><br />",
		result[1], result[2], result[1], result[2])
	log.Debug(button)
	return button
}

func init() {
	HelpString = FormatHelpString(commands)

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

// takes a gumble.TextMessageEvent and cmd
// runs the cmd if it exits and return it's output
// else return a canned response
func HandleCmd(args []string) (string, bool) {
	for _, cmd := range commands {
		if cmd.Name() == args[1] {
			return cmd.Run(cmd, args), cmd.PublicResponse
		}
	}
	return CommandNotFound(args[1]), false
}

//parse steam connect strings and provide a html button to the channel
func HandleMessage(e *gumble.TextMessageEvent, config *configuration.ButlerConfiguration) {
	//check for steam connect cmds
	result := Steamconnect.FindStringSubmatch(e.Message)
	if result != nil {
		e.Client.Self.Channel.Send(FormatSteamConnect(result), config.Bot.RecursiveChannelMessages)
	} else {
		//check for bot commands
		result = ChatCommand.FindStringSubmatch(e.Message)
		if result != nil {
			if result[1] == "help" {
				//special case to avoid  initialization loop
				e.Sender.Send(Help(result))
			} else {
				log.Infof("User %s (ID:%d) called '%s'", e.Sender.Name, e.Sender.UserID, result[0])
				responce, PublicResponse := HandleCmd(result)
				if PublicResponse { //send to channel
					e.Client.Self.Channel.Send(responce, config.Bot.RecursiveChannelMessages)
				} else { //send to usr
					e.Sender.Send(responce)
				}
			}
		}
	}
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
			HandleMessage(e, &config)
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
