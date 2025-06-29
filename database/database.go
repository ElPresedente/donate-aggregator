package database

import (
	_ "modernc.org/sqlite"
)

var CredentialsDB CredentialsDatabase

func InitDataBases() {
	CredentialsDB.Init()
}

func CloseDataBases() {
	CredentialsDB.db.Close()
}
