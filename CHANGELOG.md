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

## Sun 28 May 2017 04:13:00 PM MDT Version: 0.0.20
Wrote the bulk of the client's logic, registering, logging in and pinging
Added -register flag to client

## Sun 28 May 2017 06:18:56 PM MDT Version: 0.0.21
Moved cmd/{client,server}/version to root of directory, making botch packages use the same version package
Updated build script to be less redundant
Updated setup script to work with new build script

## Sun 28 May 2017 11:32:30 PM MDT Version: 0.0.22
Refactored routes.go
Moved/renamed setDefaultResponseHeaders to utils/SetResponseHeaders
Added scripts/test directory, moved test.sh from scripts/build to scripts/test
Added scripts/test/mock_client.sh

## Sun 28 May 2017 11:56:09 PM MDT Version: 0.0.23
Migrate to using the logrus package in the server binary
Added debug_db configuration option
Added LogDBError

## Mon 29 May 2017 02:08:41 PM MDT Version: 0.0.24
Added ErrAlreadyOnline, ErrAlreadyOffline and ErrGettingStatus to cmd/server/routes.go

## Mon 29 May 2017 02:12:40 PM MDT Version: 0.0.25
Added APIResult structure that holds the local error, api error and response json from an api call
Make all interfaces to the api return APIResult structure in connection.go

## Mon 29 May 2017 04:49:34 PM MDT Version: 0.0.26
Rename to go-agent

## Mon 29 May 2017 05:18:39 PM MDT Version: 0.0.27
Added ping_interval to client config

## Mon 29 May 2017 06:16:44 PM MDT Version: 0.0.28
Use nested configuration for database and general server configuration
Renamed etc/config.toml to etc/server_config.toml

## Mon 29 May 2017 06:33:42 PM MDT Version: 0.0.29
Added [server] section to etc/server_config.toml for general server configuration options

Added register_auth_count to server config section. This option dictates how many registration authorizations 
the server should generate if there are none in the database at server boot up

Added register_auth_expire to server config section. This option dictates how long a client has to use a registration
string before it expires and can no longer be used to register a device

## Mon 29 May 2017 07:05:03 PM MDT Version: 0.0.30
Added logging to systemd journal
Added systemd_logging to etc/server_config.toml to enable/disable systemd journal logging

## Thu 01 Jun 2017 04:53:25 PM MDT Version: 0.0.31
Added ssl_key_path, ssl_cert_path and listen_port to etc/server_config.toml

## Fri 02 Jun 2017 12:45:06 PM MDT Version: 0.0.32
Added scripts/build/gen_tls.sh to generate key/certificate for agent-server
Refactored cmd/client/main.go to be a little more modular
Updated etc/client_config.toml and etc/server_config.toml

## Fri 09 Jun 2017 07:48:59 PM MDT Version: 0.0.33
Removed some dead code in cmd/server/utils

## Sat 10 Jun 2017 04:56:02 PM MDT Version: 0.0.34
Make sure http response bodys get closed in all api endpoint handler

## Sat 10 Jun 2017 05:05:42 PM MDT Version: 0.0.35
Added timetrack to all api endpoints
Added timetrack_api option to server configuration
