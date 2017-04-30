package db

import (
	//	"database/sql"
	//	"encoding/json"
	"time"

	_ "github.com/ziutek/mymysql/godrv"
)

type LocationEntry struct {
	Ssid      string `json:"ssid"`
	Addr      string `json:"addr"`
	LoggedIn  bool   `json:"logged_in"`
	LoginName string `json:"login_name"`
}

type ClientError struct {
	Str       string
	Timestamp time.Time
	Fatal     bool
}

type Device struct {
	UUID        string           `json:"uuid"`
	Hostname    string           `json:"hostname"`
	Online      bool             `json:"online"`
	LastSeen    time.Time        `json:"last_seen"`
	LocationLog []*LocationEntry `json:"location_log"`
}

func HandleRegister(data []byte) {}

func HandleLogin(data []byte) {}

func HandleLogoff(data []byte) {}

func HandlePing(data []byte) {}

func HandleError(data []byte) {}
