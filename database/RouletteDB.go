package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type RouletteDatabase struct {
	db *sql.DB
}

type RouletteItem struct {
	ID      int    `json:"id"`
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
}

type RouletteGroupWithItems struct {
	Title      string   `json:"title"`      // name
	Items      []string `json:"items"`      // names из связанных RouletteItem
	Percentage float64  `json:"percentage"` // chance
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

//Получение всех RouletteItem по id RouletteGroup
func (c *RouletteDatabase) GetItemsByGroupID(groupID int) ([]RouletteItem, error) {
	rows, err := c.db.Query(`SELECT id, group_id, name FROM RouletteItem WHERE group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []RouletteItem
	for rows.Next() {
		var item RouletteItem
		if err := rows.Scan(&item.ID, &item.GroupID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (c *RouletteDatabase) AddItemToGroup(groupID int, itemName string) error {
	_, err := c.db.Exec(`INSERT INTO RouletteItem (group_id, name) VALUES (?, ?)`, groupID, itemName)
	return err
}

func (c *RouletteDatabase) GetGroupWithItemsByID(groupID int) (*RouletteGroupWithItems, error) {
	query := `
		SELECT g.name AS group_name, g.chance, i.name AS item_name
		FROM RouletteGroup g
		LEFT JOIN RouletteItem i ON g.id = i.group_id
		WHERE g.id = ?
	`

	rows, err := c.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *RouletteGroupWithItems

	for rows.Next() {
		var groupName string
		var chance float64
		var itemName sql.NullString

		err := rows.Scan(&groupName, &chance, &itemName)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = &RouletteGroupWithItems{
				Title:      groupName,
				Percentage: chance,
				Items:      []string{},
			}
		}

		if itemName.Valid {
			result.Items = append(result.Items, itemName.String)
		}
	}

	if result == nil {
		return nil, fmt.Errorf("Группа с id %d не найдена", groupID)
	}

	return result, nil
}