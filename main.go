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
	expire, err := time.ParseDuration("30s")
	if err != nil {
		panic(err)
	}
	registerAuth := auth.NewRegisterAuth(expire)
	log.Println("Generated registration authorization string:", registerAuth.Str)
	log.Println("Starting tracker server...")
	if err := db.Init(); err != nil {
		panic(err)
	}
	log.Printf("Initialized database connection %s...", options.Config.Addr)
	routes.RegisterHandles()
}

func main() {
	routes.Launch()
}
