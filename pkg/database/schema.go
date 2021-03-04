package database

// SQLite schema for logging data in rainbase

const (
	//nolint
	sqlschema = `
PRAGMA foreign_keys = ON; --required 
BEGIN TRANSACTION;
DROP TABLE IF EXISTS mappings;
CREATE TABLE mappings (
	id INTEGER PRIMARY KEY,
	longname TEXT
);

INSERT INTO mappings (id, longname) 
VALUES
	(0, "rain event"),
	(1, "temperature"),
	(2, "soft reset event"),
	(3, "hard reset event"),
	(4, "pause"),
	(5, "unpause"),
	(6, NULL),
	(7, NULL)
;

DROP TABLE IF EXISTS log;
CREATE TABLE log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tag INTEGER NOT NULL,
	value INTEGER NOT NULL,
	timestamp TEXT NOT NULL, --created by go
	FOREIGN KEY (tag) REFERENCES mappings(id)
);

INSERT INTO log (tag, value, timestamp)
VALUES (7,1,"timestamp");

COMMIT;
`
)
