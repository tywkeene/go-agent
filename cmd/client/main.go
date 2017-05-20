package main

import (
	"fmt"
	"os"

	"github.com/tywkeene/go-tracker/cmd/client/connection"
	"github.com/tywkeene/go-tracker/cmd/client/options"
)

func init() {
	options.ParseFlags()
}

func main() {
	c := connection.NewConnection(options.Config.ServerAddr)
	auth := options.ReadAuthFile()
	if auth.UUID == "" {
		if err := c.Register(auth.Auth); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
