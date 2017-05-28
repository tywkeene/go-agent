package options

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/tywkeene/go-tracker/cmd/server/version"
	"os"
)

type Configuration struct {
	Addr string `toml:"address"`
	Name string `toml:"name"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

var Config Configuration

func ReadConfig() {
	configFile := flag.String("config", "", "configuration file to parse")
	printVersion := flag.Bool("version", false, "print version information")
	flag.Parse()
	if *printVersion == true {
		version.Print()
		os.Exit(0)
	}
	if *configFile == "" {
		panic("Need config file")
	}
	if _, err := toml.DecodeFile(*configFile, &Config); err != nil {
		panic(err)
	}
	if Config.Addr == "" || Config.Name == "" ||
		Config.User == "" || Config.Pass == "" {
		panic("Invalid database configuration")
	}
}