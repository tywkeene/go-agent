package main

import (
	"github.com/tywkeene/go-tracker/cmd/server/auth"
	"github.com/tywkeene/go-tracker/cmd/server/db"
	"github.com/tywkeene/go-tracker/cmd/server/options"
	"github.com/tywkeene/go-tracker/cmd/server/routes"
	"log"
	"time"
)

func init() {
	options.ReadConfig()
	log.Println("Starting tracker server...")
	if err := db.Init(); err != nil {
		panic(err)
	}
	log.Printf("Initialized database connection %s...", options.Config.Addr)

	expire, err := time.ParseDuration("24h")
	if err != nil {
		panic(err)
	}
	if err := auth.Init(10, expire); err != nil {
		panic(err)
	}
	routes.RegisterHandles()
}

func main() {
	routes.Launch()
}
