package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tywkeene/go-tracker/cmd/client/connection"
	"github.com/tywkeene/go-tracker/cmd/client/options"
)

func init() {
	options.ParseFlags()
}

func main() {
	c := connection.NewConnection(options.Config.ServerAddr)
	auth := options.ReadAuthFile()
	if options.Config.NewAuth != "" {
		if err := c.Register(options.Config.NewAuth); err != nil {
			panic(err)
		}
		newAuth := &options.Authorization{
			UUID:    c.Device.UUID,
			AuthStr: c.Device.AuthStr,
		}
		options.WriteAuthFile(newAuth)
		os.Exit(0)
	}
	// Status Check to see if our uuid still exists on the server
	// If not we attempt to register with the server
	if auth.UUID != "" {
		fmt.Println("Already have a UUID, checking registration status")
		registered, err := c.GetStatus(auth)
		if err != nil {
			fmt.Println(err)
		}
		if registered == false {
			fmt.Println(`Device not registered, you can register by running:\n\t
			tracker-client -config <config file> -register <auth string>`)
			os.Exit(-1)
		} else {
			// Login
			fmt.Println("Device registered with server, logging in")
			if err := c.Login(); err != nil {
				panic(err)
			}
		}
	}

	// Ping loop
	for {
		time.Sleep(5 * time.Second)
		if err := c.Ping(); err != nil {
			panic(err)
		}
	}
}
