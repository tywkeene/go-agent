package options

import (
	"flag"
	"github.com/BurntSushi/toml"
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
	flag.Parse()
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
