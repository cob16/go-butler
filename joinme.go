package main

import (
	"fmt"
	"layeh.com/gumble/gumble"
)

var joinme = &Command{
	Run:            joinUser,
	PublicResponse: true,
	UsageLine:      "joinme",
	Short:          "Connects to channel that user is currenly in",
	Long: `
Connects to the channel of the user that called the command.

Will nottify both the source and destnation channel
`,
}

func joinUser(cmd *Command, args []string, event *gumble.TextMessageEvent) string {

	if event.Sender.Channel == event.Client.Self.Channel {
		return "Sir im am already in the room!"
	}

	var excuse string = fmt.Sprintf("'%s' is asking for my services in '%s'. Please Excuse me",
		event.Sender.Name,
		event.Sender.Channel.Name,
		)
	log.Info(excuse)
	event.Client.Self.Channel.Send(excuse, false)


	event.Client.Self.Move(event.Sender.Channel)

	return fmt.Sprintf("Moved to '%s' from '%s' by'%s'",
		event.Sender.Channel.Name,
		event.Client.Self.Channel.Name,
		event.Sender.Name,)
}