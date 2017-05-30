package main

import (
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/tywkeene/go-agent/cmd/server/auth"
	"github.com/tywkeene/go-agent/cmd/server/db"
	"github.com/tywkeene/go-agent/cmd/server/options"
	"github.com/tywkeene/go-agent/cmd/server/routes"
)

func init() {
	options.ReadConfig()
	log.Infof("Starting go-agent server...")
	if err := db.Init(); err != nil {
		panic(err)
	}
	dbConfig := options.Config.Database
	serverConfig := options.Config.Server

	log.Infof("Initialized database connection %s...", dbConfig.Addr)

	expire, err := time.ParseDuration(serverConfig.RegisterAuthExpire)
	if err != nil {
		panic(err)
	}
	if err := auth.Init(serverConfig.RegisterAuthCount, expire); err != nil {
		panic(err)
	}
	routes.RegisterHandles()
}

func main() {
	routes.Launch()
}
