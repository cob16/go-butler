package main

import (
	"github.com/njdart/go-butler/steamgauge"
	"github.com/layeh/gumble/gumble"
)

var status = &Command{
	Run:            SteamStatus,
	PublicResponse: true,
	UsageLine:      "status [tf2|dota|cs:go]",
	Short:          "shows info of online steam services",
	Long: `
Shows the status of steam services using the steamgaug.es
a seconds arg of tf2 dota or cs:go can be used to filter
infomation for that partular game
`,
}

func SteamStatus(cmd *Command, args []string, event *gumble.TextMessageEvent) string {
	status, err := steamgauge.GetSteamStatus()
	if err != nil {
		log.Panic(err)
	}
	if args[2] != "" {
		switch args[2] {
		case "tf2":
			return status.GetStatusTF2()
		case "csgo":
			return status.GetStatusCSGO()
		case "dota":
			return status.GetStatusDOTA2()
		case "steam":
			return status.GetStatusSteam()
		default:
			return CommandNotFound(args[0])
		}
	}
	return status.GetStatus()
}
