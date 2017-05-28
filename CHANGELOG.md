## Thu 04 May 2017 09:55:36 PM MDT Version: 0.0.1
Registering a device works. Added check for duplicate hostname

## Thu 04 May 2017 10:31:44 PM MDT Version: 0.0.2
Fixed DeviceByUuidStmt in db/db.go

## Thu 04 May 2017 11:12:51 PM MDT Version: 0.0.3
Removed stale code in main.go

## Thu 04 May 2017 11:39:06 PM MDT Version: 0.0.4
Added database configuration via toml file.
Added relevent file options/options.go
Added example config.toml
Updated .gitignore

## Thu 04 May 2017 11:43:17 PM MDT Version: 0.0.5
Added bind_addr option to config.toml

## Fri 05 May 2017 05:19:55 PM MDT Version: 0.0.6
Make config file just for the database connection

## Sun 07 May 2017 12:34:19 PM MDT Version: 0.0.7
Added time sensitive one-time authorization for devices
Updated schema, adding auth_string to tracker.devices
Fail if any configuration options are nil. (hack for now, fix later)

## Tue 09 May 2017 10:40:58 PM MDT Version: 0.0.8
Put registration authorizations strings in the database

## Tue 09 May 2017 11:08:35 PM MDT Version: 0.0.9
Added address (ip address of device) to device database table
Refactored auth string validation to check for valid auth string, expiration and use

## Wed 10 May 2017 12:08:41 PM MDT Version: 0.0.10
Changed default registration auth expiration to 24h
Strip port from device ip addr when registering a device

## Thu 11 May 2017 07:42:09 PM MDT Version: 0.0.11
Make errors returned from db more standardized.

## Wed 17 May 2017 12:26:40 AM MDT Version: 0.0.12
Implemented login and logoff endpoints

## Thu 18 May 2017 09:27:45 AM MDT Version: 0.0.13
Added forgotten primary key to error_reports table in schema.sql

## Thu 18 May 2017 10:50:11 AM MDT Version: 0.0.14
Implemented /ping endpoint. Added SetOnlineStatus function.
Devices will come back online if server recieves a ping from that device

## Thu 18 May 2017 12:51:19 PM MDT Version: 0.0.15
Cleaned up db/db.go.
Removed unused error variables
Removed authorizeDeviceHostName() and authorizeDeviceUUID()

## Fri 19 May 2017 03:48:32 PM MDT Version: 0.0.16
Restructure entire project to allow for building server/client binaries separately
Added -version flag to server binary

## Fri 19 May 2017 08:35:54 PM MDT Version: 0.0.17
Start work on tracker-client. Registration works

## Sun 28 May 2017 01:32:03 PM MDT Version: 0.0.18
Remove unused table in database for now
Update routes.go and db.go accordingly

## Sun 28 May 2017 02:09:33 PM MDT Version: 0.0.19
Added /status api endpoint, which merely checks for the existence of a device
by hostname, uuid and authorization. Should be the first thing any client does
before attempting to register or login with a server.
Updated db.go/AuthorizeDevice() to just return ErrAUnauthorizedDevice
