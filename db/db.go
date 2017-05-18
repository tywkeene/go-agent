package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

// API Errors

var ErrHostnameNotAuthorized = fmt.Errorf("no device with that hostname is registered with this server")
var ErrUUIDNotAuthorized = fmt.Errorf("no device with that UUID is registered with this server")

var ErrAuthUsed = fmt.Errorf("device authorization already used")
var ErrAuthExpired = fmt.Errorf("device authorization expired")
var ErrAuthStringInvalid = fmt.Errorf("invalid device authorization string")

var ErrDeviceWithHostnameExists = fmt.Errorf("a device with that hostname already registered with this server")
var ErrDeviceWithUUIDExists = fmt.Errorf("a device with that UUID already registered with this server")
var ErrLoginUnauthorized = fmt.Errorf("unauthorized login: device hostname, uuid or auth string unrecognized")
var ErrUnauthorizedDevice = fmt.Errorf("unknown or unauthorized device")

var ErrDatabaseError = fmt.Errorf("internal database error")

const RegisterStmt = "INSERT INTO devices SET uuid=?,address=?,auth_string=?,hostname=?,online=?;"
const DeviceByHostStmt = "SELECT hostname FROM devices WHERE hostname=?;"
const DeviceByUUIDStmt = "SELECT uuid FROM devices WHERE uuid=?;"
const AuthorizeDeviceStmt = "SELECT hostname,uuid,auth_string FROM devices WHERE hostname=? AND uuid=? AND auth_string=?;"
const SetOnlineStatusStmt = "UPDATE devices SET online=? WHERE hostname=? AND uuid=? AND auth_string=?;"
const GetDeviceStatusStmt = "SELECT online FROM devices WHERE uuid=?;"
const PingStmt = "UPDATE devices SET last_seen=? WHERE hostname=? AND uuid=? AND auth_string=?;"

const RegisterAuthCount = "SELECT COUNT(*) FROM register_auths;"
const InsertRegisterAuthStmt = "INSERT INTO register_auths SET auth_string=?,used=?,timestamp=?,expire_timestamp=?;"
const ValidateRegisterAuthStmt = "SELECT auth_string,used,timestamp,expire_timestamp FROM register_auths WHERE auth_string=?;"
const SetRegisterAuthUsedStmt = "UPDATE register_auths SET used=? WHERE auth_string=? ;"

func GetRegisterAuthCount() (int, error) {
	rows, err := DBConnection.Query(RegisterAuthCount)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return 0, ErrDatabaseError
	}
	rows.Next()
	var rowCount int
	err = rows.Scan(&rowCount)
	if err != nil {
		log.Println(err)
		return 0, ErrDatabaseError
	}
	return rowCount, nil
}

func InsertRegisterAuth(str string, used bool, timestamp int64, expire int64) error {
	stmt, err := DBConnection.Prepare(InsertRegisterAuthStmt)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	_, err = stmt.Exec(str, used, timestamp, expire)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
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
		log.Println(err)
		return false, ErrDatabaseError
	}

	if authStr != str {
		return false, ErrAuthStringInvalid
	} else if used == true {
		return false, ErrAuthUsed
	} else if expireTimestamp < time.Now().Unix() {
		return false, ErrAuthExpired
	}
	return true, nil
}

func SetAuthUsed(authStr string, used bool) error {
	stmt, err := DBConnection.Prepare(SetRegisterAuthUsedStmt)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
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
		log.Println(err)
		return false, ErrDatabaseError
	}
	return true, nil
}

func authorizeDeviceHostName(device *Device) error {
	exists, err := RowExists(DeviceByHostStmt, device.Hostname)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	if exists == false {
		return ErrHostnameNotAuthorized
	}
	return nil
}

func authorizeDeviceUUID(device *Device) error {
	exists, err := RowExists(DeviceByUUIDStmt, device.UUID)
	if exists == false {
		return ErrUUIDNotAuthorized
	}
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	return nil
}

func AuthorizeDevice(device *Device) (bool, error) {
	var hostname string
	var uuid string
	var auth string
	err := DBConnection.QueryRow(AuthorizeDeviceStmt,
		device.Hostname, device.UUID, device.AuthStr).Scan(&hostname, &uuid, &auth)
	if err == sql.ErrNoRows {
		return false, ErrDatabaseError
	} else if err != nil {
		log.Println(err)
		return false, ErrDatabaseError
	}
	if device.Hostname == hostname &&
		device.UUID == uuid &&
		device.AuthStr == auth {
		return true, nil
	}
	return false, nil
}

func HandleRegister(device *Device) error {
	exists, err := RowExists(DeviceByHostStmt, device.Hostname)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	if exists == true {
		return ErrDeviceWithHostnameExists
	}
	stmt, err := DBConnection.Prepare(RegisterStmt)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	_, err = stmt.Exec(device.UUID, device.Address, device.AuthStr, device.Hostname, device.Online)
	return err
}

func IsDeviceOnline(device *Device) (bool, error) {
	var online bool
	err := DBConnection.QueryRow(GetDeviceStatusStmt, device.UUID).Scan(&online)
	if err == sql.ErrNoRows || err != nil {
		log.Println(err)
		return false, ErrDatabaseError
	}
	return online, nil
}

func SetDeviceOnlineStatus(device *Device, online bool) error {
	stmt, err := DBConnection.Prepare(SetOnlineStatusStmt)
	if err != nil {
		return nil
	}
	_, err = stmt.Exec(online, device.Hostname, device.UUID, device.AuthStr)
	return err
}

func HandleLogin(device *Device) error {
	auth, err := AuthorizeDevice(device)
	if err != nil || auth == false {
		return err
	}
	if err := SetDeviceOnlineStatus(device, true); err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	return err
}

func HandleLogoff(device *Device) error {
	auth, err := AuthorizeDevice(device)
	if err != nil || auth == false {
		return err
	}
	if err := SetDeviceOnlineStatus(device, false); err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	return nil
}

func HandlePing(device *Device) error {
	auth, err := AuthorizeDevice(device)
	if err != nil || auth == false {
		return err
	}
	var timestamp = time.Now().Format(time.RFC3339)
	stmt, err := DBConnection.Prepare(PingStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(timestamp, device.Hostname, device.UUID, device.AuthStr)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}

	online, err := IsDeviceOnline(device)
	if err != nil {
		log.Println(err)
		return ErrDatabaseError
	}
	if online == false {
		stmt, err = DBConnection.Prepare(SetOnlineStatusStmt)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(true, device.Hostname, device.UUID, device.AuthStr)
		if err != nil {
			log.Println(err)
			return ErrDatabaseError
		}
	}
	return nil
}

func HandleError(data []byte) {}

func Init() error {
	var err error
	dbOptions := options.Config
	DBConnection, err = sql.Open("mysql", dbOptions.User+":"+dbOptions.Pass+"@tcp("+dbOptions.Addr+")/"+dbOptions.Name)
	return err
}
