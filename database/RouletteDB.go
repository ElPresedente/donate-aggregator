package database

import (
	"database/sql"
	//"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type RouletteDatabase struct {
	db *sql.DB
}

func (c *RouletteDatabase) Init() {
	var err error
	c.db, err = sql.Open("sqlite", "./RouletteDB.db")
	if err != nil {
		log.Printf("❌ Ошибка подключения к базе RouletteDB: %s", err)
	}

	if err = c.db.Ping(); err != nil {
		log.Fatal("❌Ошибка пинга БД:", err)
	}

	//Это высрал чатгпт, он сказал включать внешние ключи
	_, err = c.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
	log.Fatal("Ошибка включения foreign_keys:", err)
	}

	//ХЗ надо ли разбивать на подзапросы, чатгпт сказал забить болт

	//КОРОЧЕ ДОПИЛИТЬ ЧТОБЫ ИНСЕРТЫ ДЕЛАЛИСЬ ОДИН РАЗ ПРИ СОЗДАНИИ
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS RouletteGroup (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        chance REAL NOT NULL
    );
	
	CREATE TABLE IF NOT EXISTS RouletteItem (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
        name TEXT NOT NULL,
		FOREIGN KEY (group_id) REFERENCES RouletteGroup(id)
    );
	`

	_, err = c.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в RouletteDB: %s", err)
	}
	RouletteDB.seedDefaultGroups()
}

func (c *RouletteDatabase) seedDefaultGroups() {
    var count int
    err := c.db.QueryRow(`SELECT COUNT(*) FROM RouletteGroup`).Scan(&count)
    if err != nil {
        log.Printf("❌ Ошибка при проверке таблицы RouletteGroup: %s", err)
        return
    }

    if count > 0 {
        return // Данные уже есть — ничего не делаем
    }

    insertDefaults := `
        INSERT INTO RouletteGroup (name, chance) VALUES ('Обычные', 50);
        INSERT INTO RouletteGroup (name, chance) VALUES ('Необычные', 25);
        INSERT INTO RouletteGroup (name, chance) VALUES ('Редкие', 16);
        INSERT INTO RouletteGroup (name, chance) VALUES ('Эпические', 7);
        INSERT INTO RouletteGroup (name, chance) VALUES ('Легендарные', 1.5);
        INSERT INTO RouletteGroup (name, chance) VALUES ('Артифакты', 0.5);
    `
    _, err = c.db.Exec(insertDefaults)
    if err != nil {
        log.Printf("❌ Ошибка при вставке дефолтных данных в RouletteGroup: %s", err)
    } else {
        log.Println("✅ Дефолтные группы успешно добавлены.")
    }
}