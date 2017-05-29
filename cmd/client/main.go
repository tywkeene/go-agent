package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tywkeene/go-agent/cmd/client/connection"
	"github.com/tywkeene/go-agent/cmd/client/options"

	"github.com/tywkeene/go-agent/cmd/server/routes"
)

func init() {
	options.ParseFlags()
}

func main() {
	c := connection.NewConnection(options.Config.ServerAddr)
	auth := options.ReadAuthFile()
	if options.Config.NewAuth != "" {
		if result := c.Register(options.Config.NewAuth); result.Ok() == false {
			result.PrintErrors()
			os.Exit(-1)
		}
		newAuth := &options.Authorization{
			UUID:    c.Device.UUID,
			AuthStr: c.Device.AuthStr,
		}
		options.WriteAuthFile(newAuth)
		fmt.Println("Successfully registered with", c.Address)
	} else if auth.UUID != "" {
		// Status Check to see if our uuid still exists on the server
		// If not we attempt to register with the server
		fmt.Println("Already have a UUID, checking registration status...")
		registered, err := c.GetStatus(auth)
		if err != nil {
			fmt.Println(err)
		}
		if registered == false {
			fmt.Println(`Device not registered, you can register by running:\n\t
			agent-client -config <config file> -register <auth string>`)
			os.Exit(-1)
		} else {
			// Login
			fmt.Println("Device registered with server, logging in...")
			if result := c.Login(); result.Ok() == false {
				if result.APIErr.ErrorMessage != routes.ErrAlreadyOnline.Error() {
					result.PrintErrors()
				}
			}
		}
	}
	fmt.Println("Successfully logged in...")
	// Ping loop
	for {
		time.Sleep(5 * time.Second)
		if result := c.Ping(); result.Ok() == false {
			result.PrintErrors()
		}
	}
}
