package options

import (
	"flag"
	"github.com/BurntSushi/toml"
)

type DB struct {
	Addr string `toml:"address"`
	Name string `toml:"name"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type Configuration struct {
	Database *DB `toml:"database"`
}

var Config Configuration

func ReadConfig() {
	configFile := flag.String("config", "config.toml", "configuration file to parse")
	flag.Parse()
	if *configFile == "" {
		panic("Need config file")
	}

	if _, err := toml.DecodeFile(*configFile, &Config); err != nil {
		panic(err)
	}
}
