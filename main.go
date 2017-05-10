package main

import (
	"github.com/tywkeene/go-tracker/auth"
	"github.com/tywkeene/go-tracker/db"
	"github.com/tywkeene/go-tracker/options"
	"github.com/tywkeene/go-tracker/routes"
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
