package main

import (
  "fmt"
  "github.com/layeh/gumble/gumble"
  "github.com/layeh/gumble/gumbleutil"
  "crypto/tls"
  "net"
)

func main() {
  config := gumble.NewConfig()
  config.Username = "gumble-test"

  var tlsConfig tls.Config
  tlsConfig.InsecureSkipVerify = true

  keepAlive := make(chan bool)

  config.Attach(gumbleutil.Listener{
    TextMessage: func(e *gumble.TextMessageEvent) {
        fmt.Printf("Received text message: %s\n", e.Message)
    },
  })

  config.Attach(gumbleutil.Listener{
    UserChange: func(e *gumble.UserChangeEvent) {
      if e.Type.Has(gumble.UserChangeConnected) {
        e.User.Send("Welcome to the server, " + e.User.Name + "!")
      }
    },
    Disconnect: func(e *gumble.DisconnectEvent) {
      keepAlive <- true
    },
  })

  client, err := gumble.DialWithDialer(new(net.Dialer), "server.com:64738", config, &tlsConfig)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }

  for _, user := range client.Users {
      fmt.Println(user)
  }

  <-keepAlive

}
