package options

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/tywkeene/go-tracker/cmd/server/version"

	"github.com/BurntSushi/toml"
)

type Authorization struct {
	UUID    string `toml:"uuid"`
	AuthStr string `toml:"auth_string"`
}

type Configuration struct {
	// Authorization string to use if the -register flag was passed
	NewAuth    string
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

func WriteAuthFile(auth *Authorization) error {
	outfile, err := os.Create(Config.AuthFile)
	if err != nil {
		return err
	}
	defer outfile.Close()

	buff := new(bytes.Buffer)
	if err := toml.NewEncoder(buff).Encode(auth); err != nil {
		return err
	}

	_, err = outfile.WriteString(buff.String())
	return err
}

func ParseFlags() {
	configFile := flag.String("config", "", "configuration file to parse")
	printVersion := flag.Bool("version", false, "print version information")
	register := flag.String("register", "", "register with the configured server using an authorization string")
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
	if *register != "" {
		Config.NewAuth = *register
	}
}
