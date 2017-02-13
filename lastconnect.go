package main

import (
	"fmt"
	"layeh.com/gumble/gumble"
)

var connect = &Command{
	Run:            Lastconnect,
	PublicResponse: true,
	UsageLine:      "connect [delete]",
	Short:          "Shows last connect info pasted into mumble",
	Long: `
Shows last connect info last pasted into mumbble, you can allso
use the "delete" arg to remove this from memory.

note: This will replace the connect string with your user ID!
`,
}

var lastconnect string = "There is not yet a connect!"
var lastconnect_raw string = "There is not yet a connect!"

func Lastconnect(cmd *Command, args []string, event *gumble.TextMessageEvent) string {
	sender := event.Sender
	if args[2] != "" {
		switch args[2] {
		case "delete":
			lastconnect = fmt.Sprintf("The last connect was deleted by '%s' ID: %d ", sender.Name, sender.UserID)
			lastconnect_raw = lastconnect
			return lastconnect
		case "raw":
			return lastconnect_raw
		default:
			return CommandNotFound(args[0])
		}
	}
	return lastconnect
}