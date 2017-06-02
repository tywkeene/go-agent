package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tywkeene/go-agent/cmd/client/connection"
	"github.com/tywkeene/go-agent/cmd/client/options"

	"github.com/tywkeene/go-agent/cmd/server/routes"
)

var serverConn *connection.Connection
var clientAuth *options.Authorization

func init() {
	options.ParseFlags()
	serverConn = connection.NewConnection(options.Config.ServerAddr)
	clientAuth = options.ReadAuthFile()
}

func doAuthorization() {
	// A registration authorization string was supplied on the commandline
	// Try to use it to register with the server.
	if result := serverConn.Register(options.Config.NewAuth); result.Ok() == false {
		result.PrintErrors()
		os.Exit(-1)
	}
	newAuth := &options.Authorization{
		UUID:    serverConn.Device.UUID,
		AuthStr: serverConn.Device.AuthStr,
	}
	options.WriteAuthFile(newAuth)
	fmt.Println("Successfully registered with", serverConn.Address)
}

func doLogin() {
	// Status Check to see if our uuid still exists on the server
	// If not we attempt to register with the server
	fmt.Println("Already have a UUID, checking registration status...")
	registered, err := serverConn.GetStatus(clientAuth)
	if err != nil {
		fmt.Println(err)
	}
	if registered == false {
		fmt.Printf("Device not registered, you can register by running:\n\t" +
			"$ agent-client -config <config file> -register <auth string>\n" +
			"to register this device with the server, and then run again\n")
		os.Exit(-1)
	} else {
		// Login
		fmt.Println("Device registered with server, logging in...")
		if result := serverConn.Login(); result.Ok() == false {
			if result.APIErr.ErrorMessage != routes.ErrAlreadyOnline.Error() {
				result.PrintErrors()
				os.Exit(-1)
			}
		}
	}
}

func main() {
	if options.Config.NewAuth != "" {
		doAuthorization()
	} else if clientAuth.UUID != "" {
		doLogin()
	}
	fmt.Println("Successfully logged in...")
	pingInterval, err := time.ParseDuration(options.Config.PingInterval)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(pingInterval)
		if result := serverConn.Ping(); result.Ok() == false {
			result.PrintErrors()
		}
	}
}
