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

type ENVVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//названия полей я задаю с фронта
//на текущий момент мы сохраняем
//	donattyToken 	- токен donatty
//	donattyUrl		- URL donatty
//	donatpayToken	- токен donatpay
//	donatpayUserId	- UserId donatpay
//	rollPrice 		- цена прокрута рулетки

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
        name TEXT PRIMARY KEY NOT NULL,
        value TEXT NOT NULL
    );`

	_, err = c.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в CredentialsDB: %s", err)
	}

	c.InitDefaultVariable()
}

// Я немножко насрал
func (c *CredentialsDatabase) checkRecordsExist(names []string) (map[string]bool, error) {
	exists := make(map[string]bool)
	for _, name := range names {
		var count int
		query := "SELECT COUNT(*) FROM EnvVariables WHERE name = ?"
		err := c.db.QueryRow(query, name).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("error checking name %s: %v", name, err)
		}
		exists[name] = count > 0
	}
	return exists, nil
}

func (c *CredentialsDatabase) InitDefaultVariable() {
	names := []string{"donattyToken", "donattyUrl", "donatpayToken", "donatpayUserId", "donatpayDomain", "rollPrice", "rollPriceIncrease", "logEnabled"}
	exists, err := c.checkRecordsExist(names)

	if err != nil {
		log.Fatal(err)
	}

	values := map[string]string{
		"donattyToken":      "",
		"donattyUrl":        "",
		"donatpayToken":     "",
		"donatpayUserId":    "",
		"donatpayDomain":    ".eu",
		"rollPrice":         "100",
		"rollPriceIncrease": "0",
		"logEnabled":        "false",
	}

	for key, value := range values {
		if !exists[key] {
			c.InsertENVValue(key, value)
		}
	}
}

func (c *CredentialsDatabase) InsertENVValue(name, value string) {
	insertQuery := "INSERT INTO EnvVariables (name, value) VALUES (?,?)"

	_, err := c.db.Exec(insertQuery, name, value)
	if err != nil {
		log.Printf("❌ Ошибка записи данных (%s:%s) в CredentialsDB: %s", name, value, err)
	}
}

func (c *CredentialsDatabase) UpdateENVValue(name, value string) {
	insertQuery := "UPDATE EnvVariables SET value = ? WHERE name = ?"

	_, err := c.db.Exec(insertQuery, value, name)
	if err != nil {
		log.Printf("❌ Ошибка записи данных (%s:%s) в CredentialsDB: %s", name, value, err)
	}
}

func (c *CredentialsDatabase) GetENVValue(name string) (string, error) {
	query := "SELECT value FROM EnvVariables WHERE name = ?"

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

func (c *CredentialsDatabase) GetAllENVValues() ([]ENVVariable, error) {
	query := "SELECT name, value FROM EnvVariables"

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envVariables []ENVVariable
	for rows.Next() {
		var variable ENVVariable
		if err := rows.Scan(&variable.Name, &variable.Value); err != nil {
			return nil, err
		}
		envVariables = append(envVariables, variable)
	}
	return envVariables, nil
}

func (c *CredentialsDatabase) CheckENVExists(name string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM EnvVariables WHERE name = ?"

	err := c.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func TestInsertGet() {
	CredentialsDB.InsertENVValue("TestName", "TestValue")
	result, _ := CredentialsDB.GetENVValue("TestName")
	log.Printf("ПРОВЕРКА ЗАПИСИ/ЧТЕНИЯ ИЗ БАЗЫ CredentialsDB: %s", result)
}
