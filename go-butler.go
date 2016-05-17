// A mumble bot based on the Gumble libary
// https://github.com/layeh/gumble/
package main

import (
  "fmt"
  "net"

  "github.com/layeh/gumble/gumble"
  "github.com/layeh/gumble/gumbleutil"
  "./configuration"
)

func main() {
  configuration, err := configuration.LoadConfiguration()
  if err != nil {
    panic(err)
  }
  tlsConfig, gumbleConfig := configuration.ExplodeConfiguration()
  log := configuration.GetLogger()
  log.Info("go-butler has sucessfully started!")

  keepAlive := make(chan bool)

  gumbleConfig.Attach(gumbleutil.Listener{
    UserChange: func(e *gumble.UserChangeEvent) {
      if e.Type.Has(gumble.UserChangeConnected) {
        e.User.Send("Welcome to the server, " + e.User.Name + "!")
      }
    },
    TextMessage: func(e *gumble.TextMessageEvent) {
      fmt.Printf("Received text message: %s\n", e.Message)
    },
    //kill the program if we are disconnected
    Disconnect: func(e *gumble.DisconnectEvent) {
      log.Info("gumble was Disconnected from the server")
      keepAlive <- true
    },
  })

  client, err := gumble.DialWithDialer(new(net.Dialer), configuration.GetUri(), &gumbleConfig, &tlsConfig)
  if err != nil {
    log.Panic(err)
  }

  for _, user := range client.Users {
    log.Info(user)
  }

  <-keepAlive
}
