package database

import (
	_ "modernc.org/sqlite"
)

var CredentialsDB CredentialsDatabase
var RouletteDB RouletteDatabase

func InitDataBases() {
	CredentialsDB.Init()
	RouletteDB.Init()
}

func CloseDataBases() {
	CredentialsDB.db.Close()
	RouletteDB.db.Close()
}
