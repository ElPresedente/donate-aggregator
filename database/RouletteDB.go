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

type RouletteGroup struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Percentage float64 `json:"chance"`
	Color      string  `json:"color"`
}

type RouletteItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RouletteGroupWithItems struct {
	Title      string   `json:"title"`      // name
	Items      []string `json:"items"`      // names из связанных RouletteItem
	Percentage float64  `json:"percentage"` // chance
}

func (rd *RouletteDatabase) Init() {
	var err error
	rd.db, err = sql.Open("sqlite", "./RouletteDB.db")
	if err != nil {
		log.Printf("❌ Ошибка подключения к базе RouletteDB: %s", err)
	}

	if err = rd.db.Ping(); err != nil {
		log.Fatal("❌Ошибка пинга БД:", err)
	}

	//Это высрал чатгпт, он сказал включать внешние ключи
	_, err = rd.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Ошибка включения foreign_keys:", err)
	}

	//ХЗ надо ли разбивать на подзапросы, чатгпт сказал забить болт

	//КОРОЧЕ ДОПИЛИТЬ ЧТОБЫ ИНСЕРТЫ ДЕЛАЛИСЬ ОДИН РАЗ ПРИ СОЗДАНИИ
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS RouletteGroup (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        chance REAL NOT NULL,
		color TEXT NOT NULL
    );
	
	CREATE TABLE IF NOT EXISTS RouletteItem (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
        name TEXT NOT NULL,
		FOREIGN KEY (group_id) REFERENCES RouletteGroup(id)
    );
	`

	_, err = rd.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в RouletteDB: %s", err)
	}
	RouletteDB.seedDefaultGroups()
}

func (rd *RouletteDatabase) seedDefaultGroups() {
	var count int
	err := rd.db.QueryRow(`SELECT COUNT(*) FROM RouletteGroup`).Scan(&count)
	if err != nil {
		log.Printf("❌ Ошибка при проверке таблицы RouletteGroup: %s", err)
		return
	}

	if count > 0 {
		return // Данные уже есть — ничего не делаем
	}

	tx, err := rd.db.Begin()
	if err != nil {

		log.Printf("❌ Ошибка начала транзакции: %s", err)
		return
	}

	groups := []struct {
		name   string
		chance float64
		color  string
	}{
		{"Обычные", 50, "white"},
		{"Необычные", 25, "rgb(55, 255, 0)"},
		{"Редкие", 16, "rgb(0, 200, 255)"},
		{"Эпические", 7, "rgb(255, 0, 251)"},
		{"Легендарные", 1.5, "rgb(245, 117, 7)"},
		{"Артефакты", 0.5, "rgb(229, 204, 128)"},
	}

	for _, g := range groups {
		_, err = tx.Exec(`INSERT INTO RouletteGroup (name, chance, color) VALUES (?, ?, ?)`, g.name, g.chance, g.color)
		if err != nil {
			tx.Rollback()
			log.Printf("❌ Ошибка вставки группы: %s", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Ошибка коммита транзакции: %s", err)
		return
	}
}

// Получение всех RouletteItem по id RouletteGroup
func (rd *RouletteDatabase) GetItemsByGroupID(groupID int) ([]RouletteItem, error) {
	rows, err := rd.db.Query(`SELECT id, name FROM RouletteItem WHERE group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []RouletteItem
	for rows.Next() {
		var item RouletteItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (rd *RouletteDatabase) AddItem(groupID int, name string) error {
	tx, err := rd.db.Begin()
	if err != nil {

		log.Printf("❌ Ошибка начала транзакции: %s", err)
		return err
	}

	_, err = tx.Exec(`INSERT INTO RouletteItem (group_id, name) VALUES (?, ?)`, groupID, name)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Ошибка вставки предмета: %s", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Ошибка коммита транзакции: %s", err)
		return err
	}
	return err
}

func (rd *RouletteDatabase) UpdateItem(id int, name string) error {
	_, err := rd.db.Exec(`UPDATE RouletteItem SET name = ? WHERE id = ?`, name, id)
	return err
}

func (rd *RouletteDatabase) DeleteItem(id int) error {
	_, err := rd.db.Exec(`DELETE FROM RouletteItem WHERE id = ?`, id)
	return err
}

func (rd *RouletteDatabase) GetRouletteGroups() ([]RouletteGroup, error) {
	rows, err := rd.db.Query(`SELECT * FROM RouletteGroup`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []RouletteGroup
	for rows.Next() {
		var group RouletteGroup
		if err := rows.Scan(&group.ID, &group.Name, &group.Percentage, &group.Color); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
