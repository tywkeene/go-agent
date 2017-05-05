package main

import (
	"github.com/tywkeene/go-tracker/db"
	"github.com/tywkeene/go-tracker/options"
	"github.com/tywkeene/go-tracker/routes"
	"log"
)

func init() {
	options.ReadConfig()
	log.Println("Starting tracker server...")
	if err := db.Init(); err != nil {
		panic(err)
	}
	log.Println("Initialized database connection...")
	routes.RegisterHandles()
}

func main() {
	routes.Launch()
}
