package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type LocationEntry struct {
	Ssid      string `json:"ssid"`
	Addr      string `json:"addr"`
	LoginName string `json:"login_name"`
}

type ClientError struct {
	Str       string    `json:"err_str"`
	Timestamp time.Time `json:"timestamp"`
	Fatal     bool      `json:"fatal"`
}

type Device struct {
	UUID        string           `json:"uuid"`
	Hostname    string           `json:"hostname"`
	Online      bool             `json:"online"`
	LastSeen    *time.Time       `json:"last_seen"`
	LocationLog []*LocationEntry `json:"location_log"`
}

var DBConnection *sql.DB

const RegisterStmt = "INSERT INTO devices SET uuid=?,hostname=?,online=?;"
const DeviceByHostStmt = "SELECT hostname FROM devices WHERE hostname=?;"
const DeviceByUUIDStmt = "SELECT uuid FROM devices WHERE uuid=?;"

func RowExists(stmt string, args ...interface{}) (bool, error) {
	var exists string
	err := DBConnection.QueryRow(stmt, args...).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func HandleRegister(uuid string, device *Device) error {
	exists, err := RowExists(DeviceByHostStmt, device.Hostname)
	if err != nil {
		return err
	}
	if exists == true {
		return fmt.Errorf("Device with that hostname already exists")
	}
	stmt, err := DBConnection.Prepare(RegisterStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(uuid, device.Hostname, true)
	return err
}

func HandleLogin(uuid string, device *Device) {
}

func HandleLogoff(data []byte) {}
func HandlePing(data []byte)   {}
func HandleError(data []byte)  {}

func Init() error {
	var err error
	DBConnection, err = sql.Open("mysql", "tracker:tracker@tcp(0.0.0.0:3306)/tracker")
	return err
}
