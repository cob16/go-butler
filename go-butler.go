package main

import (
  "fmt"
  "github.com/layeh/gumble/gumble"
  "github.com/layeh/gumble/gumbleutil"
)

func main() {
  config := gumble.NewConfig()
  config.Username = "gumble-test"

  config.Attach(gumbleutil.Listener{
    TextMessage: func(e *gumble.TextMessageEvent) {
        fmt.Printf("Received text message: %s\n", e.Message)
    },
  })

  _, err := gumble.Dial("server.com:64738", config)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}
