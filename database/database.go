package database

import (
	_ "modernc.org/sqlite"
)

var CredentialsDB CredentialsDatabase
var WidgetDB WidgetsDatabase
var LogDB LogDatabase

func InitDatabases() {
	CredentialsDB.Init()
	WidgetDB.Init()
	LogDB.Init()
}

func CloseDatabases() {
	CredentialsDB.db.Close()
	WidgetDB.db.Close()
	LogDB.db.Close()
}
