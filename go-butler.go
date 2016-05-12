// A mumble bot based on the Gumble libary
// https://github.com/layeh/gumble/
package main

import (
	"fmt"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"crypto/tls"
	"net"
	"log"
	"os"
	"io"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

//logs everything to the same file and Stdout
func InitLog(FileHandle io.Writer) {
	Handle := io.MultiWriter(FileHandle, os.Stdout)

	Info = log.New(Handle,
	"INFO: ",
	log.Ldate | log.Ltime | log.Lshortfile)

	Warning = log.New(Handle,
	"WARNING: ",
	log.Ldate | log.Ltime | log.Lshortfile)

	Error = log.New(Handle,
	"ERROR: ",
	log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	file, err := os.OpenFile("log_butler.txt", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	InitLog(file)


	Info.Println("Special Information")
	Warning.Println("There is something you need to know about")
	Error.Println("Something has failed")

	config := gumble.NewConfig()
	config.Username = "better-gumble-test"

	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true

	keepAlive := make(chan bool)

	config.Attach(gumbleutil.Listener{
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
			Info.Println("gumble was Disconnected from the server")
			keepAlive <- true
		},
	})

	client, err := gumble.DialWithDialer(new(net.Dialer), "example.com:64738", config, &tlsConfig)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, user := range client.Users {
		fmt.Println(user)
	}

	<-keepAlive
}
