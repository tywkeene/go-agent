package main

import (
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/tywkeene/go-tracker/cmd/server/auth"
	"github.com/tywkeene/go-tracker/cmd/server/db"
	"github.com/tywkeene/go-tracker/cmd/server/options"
	"github.com/tywkeene/go-tracker/cmd/server/routes"
)

func init() {
	options.ReadConfig()
	log.Infof("Starting tracker server...")
	if err := db.Init(); err != nil {
		panic(err)
	}
	log.Infof("Initialized database connection %s...", options.Config.Addr)

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
