package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"

	"github.com/tywkeene/go-tracker/options"
)

type DeviceRegister struct {
	Hostname string `json:"hostname"`
	AuthStr  string `json:"auth_string"`
}

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
	Address     string           `json:"address"`
	AuthStr     string           `json:"auth_string"`
	Hostname    string           `json:"hostname"`
	Online      bool             `json:"online"`
	LastSeen    *time.Time       `json:"last_seen"`
	LocationLog []*LocationEntry `json:"location_log"`
}

var DBConnection *sql.DB

const RegisterStmt = "INSERT INTO devices SET uuid=?,address=?,auth_string=?,hostname=?,online=?;"
const DeviceByHostStmt = "SELECT hostname FROM devices WHERE hostname=?;"
const DeviceByUUIDStmt = "SELECT uuid FROM devices WHERE uuid=?;"

const RegisterAuthCount = "SELECT COUNT(*) FROM register_auths;"
const InsertRegisterAuthStmt = "INSERT INTO register_auths SET auth_string=?,used=?,timestamp=?,expire_timestamp=?;"
const ValidateRegisterAuthStmt = "SELECT auth_string,used,timestamp,expire_timestamp FROM register_auths WHERE auth_string=?;"
const SetRegisterAuthUsedStmt = "UPDATE register_auths SET used=? WHERE auth_string=? ;"

func GetRegisterAuthCount() (int, error) {
	rows, err := DBConnection.Query(RegisterAuthCount)
	defer rows.Close()
	if err != nil {
		return 0, err
	}
	rows.Next()
	var rowCount int
	err = rows.Scan(&rowCount)
	if err != nil {
		return 0, err
	}
	return rowCount, nil
}

func InsertRegisterAuth(str string, used bool, timestamp int64, expire int64) error {
	stmt, err := DBConnection.Prepare(InsertRegisterAuthStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(str, used, timestamp, expire)
	if err != nil {
		return nil
	}
	return nil
}

func IsAuthValid(authStr string) (bool, error) {
	var str string
	var used bool
	var timestamp int64
	var expireTimestamp int64
	err := DBConnection.QueryRow(ValidateRegisterAuthStmt,
		authStr).Scan(&str, &used, &timestamp, &expireTimestamp)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if authStr != str {
		return false, fmt.Errorf("Unauthorized: Invalid auth string")
	} else if used == true {
		return false, fmt.Errorf("Unauthorized: Auth already used")
	} else if expireTimestamp < time.Now().Unix() {
		return false, fmt.Errorf("Unauthorized: Auth expired")
	}
	return true, nil
}

func SetAuthUsed(authStr string, used bool) error {
	stmt, err := DBConnection.Prepare(SetRegisterAuthUsedStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(used, authStr)
	return nil
}

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

func authorizeDeviceHostName(device *Device) error {
	exists, err := RowExists(DeviceByHostStmt, device.Hostname)
	if err != nil {
		return err
	}
	if exists == true {
		return fmt.Errorf("Device with that hostname already exists")
	}
	return nil
}

func authorizeDeviceUUID(uuid string, device *Device) error {
	exists, err := RowExists(DeviceByUUIDStmt, device.UUID)
	if err != nil {
		return err
	}
	if exists == false {
		return fmt.Errorf("Device with that UUID does not exist")
	}
	return nil
}

func HandleRegister(device *Device) error {
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
	_, err = stmt.Exec(device.UUID, device.Address, device.AuthStr, device.Hostname, device.Online)
	return err
}

func HandleLogin(uuid string, device *Device) {
}

func HandleLogoff(data []byte) {}
func HandlePing(data []byte)   {}
func HandleError(data []byte)  {}

func Init() error {
	var err error
	dbOptions := options.Config
	DBConnection, err = sql.Open("mysql", dbOptions.User+":"+dbOptions.Pass+"@tcp("+dbOptions.Addr+")/"+dbOptions.Name)
	return err
}
