package database

import (
	"database/sql"
	"go-back/l2db"

	//"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type LogDatabase struct {
	db *sql.DB
}

type RouletteLog struct {
	User string `json:"user"`
	Item string `json:"item"`
	Time string `json:"time"`
}

func (c *LogDatabase) Init() {
	var err error
	c.db, err = sql.Open("sqlite", "./LogDB.db")
	if err != nil {
		log.Printf("❌ Ошибка подключения к базе LogDatabase: %s", err)
	}

	if err = c.db.Ping(); err != nil {
		log.Fatal("❌Ошибка пинга БД:", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS RouletteLog (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user TEXT NOT NULL,
		item TEXT NOT NULL,
		time TEXT NOT NULL
		);
	`

	/*
		user: пользователь, для которого активировалась рулетка
		item: сектор, выпавший на рулетке
		time: время активации рулетки DD.MM HH.MM
	*/

	_, err = c.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в LB: %s", err)
	}
}

func (c *LogDatabase) InsertValue(user, item, time string) {
	insertQuery := "INSERT INTO RouletteLog (user, item, time) VALUES (?,?,?)"

	_, err := c.db.Exec(insertQuery, user, item, time)
	if err != nil {
		log.Printf("❌ Ошибка записи данных (%s:%s:%s) в LogDB: %s", user, item, time, err)
	}
}

func (c *LogDatabase) InsertSpins(data l2db.ResponseData) {
	for _, spin := range data.Spins {
		currentTime := time.Now().Format("02.01 15:04")
		c.InsertValue(data.User, spin.WinnerSector, currentTime)
	}
}

// Подумать над названием функции
func (c *LogDatabase) GetLastNLogs(limit int) ([]RouletteLog, error) {
	query := "SELECT user, item, time FROM RouletteLog ORDER BY id DESC LIMIT ?"

	rows, err := c.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []RouletteLog
	for rows.Next() {
		var log RouletteLog
		if err := rows.Scan(&log.User, &log.Item, &log.Time); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
