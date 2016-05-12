package main

import (
  "fmt"
  "github.com/layeh/gumble/gumble"
  "github.com/layeh/gumble/gumbleutil"

  "net"
  "github.com/njdart/go-butler/configuration"
)

func main() {
  configuration, err := configuration.LoadConfiguration()
  if err != nil {
    panic(err)
  }

  tlsConfig, gumbleConfig := configuration.ExplodeConfiguration()

  keepAlive := make(chan bool, 1)

  gumbleConfig.Attach(gumbleutil.Listener{
    TextMessage: func(e *gumble.TextMessageEvent) {
      fmt.Printf("Received text message: %s\n", e.Message)
    },
    UserChange: func(e *gumble.UserChangeEvent) {
      if e.Type.Has(gumble.UserChangeConnected) {
        e.User.Send("Welcome to the server, " + e.User.Name + "!")
      }
    },
    Connect: func(e *gumble.ConnectEvent) {

    },
  })

  client, err := gumble.DialWithDialer(new(net.Dialer), configuration.GetUri(), &gumbleConfig, &tlsConfig)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }

  for _, user := range client.Users {
      fmt.Println(user)
  }

  <-keepAlive

}
