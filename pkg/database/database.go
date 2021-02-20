// Package database stores weather data locally
// For data duplication/verification purposes, we store all packets in a sqlite database on the Pi
package database

// the main packet table, probably all we will need for such a simple project
const schema = `CREATE TABLE packet
id INT PRIMARY KEY SERIAL,
tag INT NOT NULL,
value INT NOT NULL,
timestamp TEXT /* created by go */
`
