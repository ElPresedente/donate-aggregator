package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type CredentialsDatabase struct {
	db *sql.DB
}

func (c *CredentialsDatabase) Init() {
	var err error
	c.db, err = sql.Open("sqlite", "./CredentialsDB.db")
	if err != nil {
		log.Printf("❌ Ошибка подключения к базе CredentialsDB: %s", err)
	}

	if err = c.db.Ping(); err != nil {
		log.Fatal("❌Ошибка пинга БД:", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS EnvVariables (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        value TEXT NOT NULL
    );`

	_, err = c.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в CredentialsDB: %s", err)
	}
}

func (c *CredentialsDatabase) InsertENVValue(name, value string) {
	insertQuery := "INSERT INTO EnvVariables (name, value) VALUES (?,?)"

	_, err := c.db.Exec(insertQuery, name, value)
	if err != nil {
		log.Printf("❌ Ошибка записи данных (%s:%s) в CredentialsDB: %s", name, value, err)
	}
}

func (c *CredentialsDatabase) GetENVValue(name string) (string, error) {
	query := `SELECT value FROM EnvVariables WHERE name = ?`

	var ENVValue string
	err := c.db.QueryRow(query, name).Scan(&ENVValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("❌ ENV параметра с именем '%s' не найден", name)
		}
		return "", err
	}

	return ENVValue, nil
}

func TestInsertGet() {
	CredentialsDB.InsertENVValue("TestName", "TestValue")
	result, _ := CredentialsDB.GetENVValue("TestName")
	log.Printf("ПРОВЕРКА ЗАПИСИ/ЧТЕНИЯ ИЗ БАЗЫ CredentialsDB: %s", result)
}
