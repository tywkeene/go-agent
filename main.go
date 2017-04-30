package main

import (
	"github.com/tywkeene/go-tracker/routes"
	"log"
)

func init() {
	log.Println("Starting tracker server...")
	routes.RegisterHandles()
}

func main() {
	routes.Launch()
}
