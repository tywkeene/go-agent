package options

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tywkeene/go-tracker/cmd/server/version"
	"os"
)

type Authorization struct {
	UUID string `toml:"uuid"`
	Auth string `toml:"auth_string"`
}

type Configuration struct {
	ServerAddr string `toml:"server_address"`
	AuthFile   string `toml:"auth_file"`
}

var Config Configuration

func ReadConfig(configFile string) {
	if _, err := toml.DecodeFile(configFile, &Config); err != nil {
		panic(err)
	}
	if Config.ServerAddr == "" {
		panic(fmt.Errorf("need a server address to connect"))
	}
}

func ReadAuthFile() *Authorization {
	var auth *Authorization
	if _, err := toml.DecodeFile(Config.AuthFile, &auth); err != nil {
		panic(err)
	}
	return auth
}

func ParseFlags() {
	configFile := flag.String("config", "", "configuration file to parse")
	printVersion := flag.Bool("version", false, "print version information")
	flag.Parse()
	if *printVersion == true {
		version.Print()
		os.Exit(0)
	}
	if *configFile == "" {
		panic("Need config file")
	} else {
		ReadConfig(*configFile)
	}
	if Config.AuthFile == "" {
		panic(fmt.Errorf("need path to an authorization file"))
	}
}
