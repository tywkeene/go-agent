CREATE DATABASE IF NOT EXISTS tracker;
USE mysql;
CREATE USER 'tracker'@'localhost' IDENTIFIED BY 'tracker';
GRANT ALL PRIVILEGES ON tracker . * TO 'tracker'@'localhost';
USE tracker;

CREATE TABLE IF NOT EXISTS location_entries(
    id INT NOT NULL AUTO_INCREMENT,
    ssid VARCHAR(32) NOT NULL,
    addr VARCHAR(15) NOT NULL,
    login_name VARCHAR(16) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS devices (
    id INT NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(38) NOT NULL,
    auth_string VARCHAR(16) NOT NULL,
    hostname VARCHAR(16) NOT NULL,
    online BOOLEAN NOT NULL,
    last_seen TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS error_reports (
    uuid VARCHAR(38) NOT NULL,
    string VARCHAR(64) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    fatal BOOLEAN NOT NULL
);
