package database

import (
	_ "modernc.org/sqlite"
)

var CredentialsDB CredentialsDatabase
var RouletteDB RouletteDatabase
var LogDB LogDatabase

func InitDataBases() {
	CredentialsDB.Init()
	RouletteDB.Init()
	LogDB.Init()
}

func CloseDataBases() {
	CredentialsDB.db.Close()
	RouletteDB.db.Close()
	LogDB.db.Close()
}
