package options

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tywkeene/go-agent/version"
	"os"
)

type DBConfig struct {
	Addr  string `toml:"address"`
	Name  string `toml:"name"`
	User  string `toml:"user"`
	Pass  string `toml:"pass"`
	Debug bool   `toml:"debug"`
}

type Configuration struct {
	Database DBConfig `toml:"database"`
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
	if Config.Database.Debug == true {
		fmt.Println("Will print database debug messages")
	}
	if Config.Database.Addr == "" || Config.Database.Name == "" ||
		Config.Database.User == "" || Config.Database.Pass == "" {
		panic("Invalid database configuration")
	}
}
