package database

import (
	_ "modernc.org/sqlite"
)

var CredentialsDB CredentialsDatabase
var RouletteDB RouletteDatabase
var LogDB LogDatabase

func InitDatabases() {
	CredentialsDB.Init()
	RouletteDB.Init()
	LogDB.Init()
}

func CloseDatabases() {
	CredentialsDB.db.Close()
	RouletteDB.db.Close()
	LogDB.db.Close()
}
