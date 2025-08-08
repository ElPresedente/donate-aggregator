package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type WidgetsDatabase struct {
	db *sql.DB
}

type RouletteCategory struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Percentage float64 `json:"chance"`
	Color      string  `json:"color"`
}

type RouletteSector struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RouletteSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (wd *WidgetsDatabase) Init() {
	var err error
	wd.db, err = sql.Open("sqlite", "./WidgetDB.db")
	if err != nil {
		log.Printf("❌ Ошибка подключения к базе WidgetDB: %s", err)
	}

	if err = wd.db.Ping(); err != nil {
		log.Fatal("❌Ошибка пинга БД:", err)
	}

	_, err = wd.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Ошибка включения foreign_keys:", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS RouletteSettings (
        name TEXT PRIMARY KEY NOT NULL,
        value TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS RouletteCategory (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        chance REAL NOT NULL,
		color TEXT NOT NULL
    );
	
	CREATE TABLE IF NOT EXISTS RouletteSector (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
        name TEXT NOT NULL,
		FOREIGN KEY (category_id) REFERENCES RouletteCategory(id)
    );
	`

	_, err = wd.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в WidgetDB: %s", err)
	}
	WidgetDB.initDefaultData()
}

func (wd *WidgetsDatabase) initDefaultData() {
	WidgetDB.initDefaultSettings()
	WidgetDB.initDefaultCategorys()
}

func (wd *WidgetsDatabase) checkSettingsExist(names []string) (map[string]bool, error) {
	exists := make(map[string]bool)
	for _, name := range names {
		var count int
		query := "SELECT COUNT(*) FROM RouletteSettings WHERE name = ?"
		err := wd.db.QueryRow(query, name).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("error checking name %s: %v", name, err)
		}
		exists[name] = count > 0
	}
	return exists, nil
}

func (wd *WidgetsDatabase) initDefaultSettings() {
	names := []string{"rollPrice", "rollPriceIncrease"}
	exists, err := wd.checkSettingsExist(names)

	if err != nil {
		log.Fatal(err)
	}

	values := map[string]string{
		"rollPrice":         "100",
		"rollPriceIncrease": "0",
	}

	for key, value := range values {
		if !exists[key] {
			wd.InsertRouletteSettingValue(key, value)
		}
	}
}

func (wd *WidgetsDatabase) initDefaultCategorys() {
	var count int
	err := wd.db.QueryRow(`SELECT COUNT(*) FROM RouletteCategory`).Scan(&count)
	if err != nil {
		log.Printf("❌ Ошибка при проверке таблицы RouletteCategory: %s", err)
		return
	}

	if count > 0 {
		return // Данные уже есть — ничего не делаем
	}

	tx, err := wd.db.Begin()
	if err != nil {

		log.Printf("❌ Ошибка начала транзакции: %s", err)
		return
	}

	sectors := []struct {
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

	for _, g := range sectors {
		_, err = tx.Exec(`INSERT INTO RouletteCategory (name, chance, color) VALUES (?, ?, ?)`, g.name, g.chance, g.color)
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

// Получение всех RouletteSector по id RouletteCategory
func (wd *WidgetsDatabase) GetSectorsByCategoryID(categoryID int) ([]RouletteSector, error) {
	rows, err := wd.db.Query(`SELECT id, name FROM RouletteSector WHERE category_id = ?`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Sectors []RouletteSector
	for rows.Next() {
		var Sector RouletteSector
		if err := rows.Scan(&Sector.ID, &Sector.Name); err != nil {
			return nil, err
		}
		Sectors = append(Sectors, Sector)
	}
	return Sectors, nil
}

func (wd *WidgetsDatabase) AddSector(categoryID int, name string) error {
	tx, err := wd.db.Begin()
	if err != nil {

		log.Printf("❌ Ошибка начала транзакции: %s", err)
		return err
	}

	_, err = tx.Exec(`INSERT INTO RouletteSector (category_id, name) VALUES (?, ?)`, categoryID, name)
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

func (wd *WidgetsDatabase) UpdateSector(id int, name string) error {
	_, err := wd.db.Exec(`UPDATE RouletteSector SET name = ? WHERE id = ?`, name, id)
	return err
}

func (wd *WidgetsDatabase) DeleteSector(id int) error {
	_, err := wd.db.Exec(`DELETE FROM RouletteSector WHERE id = ?`, id)
	return err
}

func (wd *WidgetsDatabase) GetRouletteCategorys() ([]RouletteCategory, error) {
	rows, err := wd.db.Query(`SELECT * FROM RouletteCategory`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectors []RouletteCategory
	for rows.Next() {
		var sector RouletteCategory
		if err := rows.Scan(&sector.ID, &sector.Name, &sector.Percentage, &sector.Color); err != nil {
			return nil, err
		}
		sectors = append(sectors, sector)
	}
	return sectors, nil
}

func (wd *WidgetsDatabase) InsertRouletteSettingValue(name, value string) {
	insertQuery := "INSERT INTO RouletteSettings (name, value) VALUES (?,?)"

	_, err := wd.db.Exec(insertQuery, name, value)
	if err != nil {
		log.Printf("❌ Ошибка записи данных (%s:%s) в InsertRouletteSettingValue: %s", name, value, err)
	}
}

func (wd *WidgetsDatabase) UpdateRouletteSettingValue(name, value string) {
	insertQuery := "UPDATE RouletteSettings SET value = ? WHERE name = ?"

	_, err := wd.db.Exec(insertQuery, value, name)
	if err != nil {
		log.Printf("❌ Ошибка записи данных (%s:%s) в UpdateRouletteSettingValue: %s", name, value, err)
	}
}

func (wd *WidgetsDatabase) GetRouletteSettingValue(name string) (string, error) {
	query := "SELECT value FROM RouletteSettings WHERE name = ?"

	var ENVValue string
	err := wd.db.QueryRow(query, name).Scan(&ENVValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("❌ Roulett Setting с именем '%s' не найден", name)
		}
		return "", err
	}

	return ENVValue, nil
}

func (wd *WidgetsDatabase) GetRouletteSettings() ([]RouletteSetting, error) {
	query := "SELECT name, value FROM RouletteSettings"

	rows, err := wd.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rouletteSettings []RouletteSetting
	for rows.Next() {
		var variable RouletteSetting
		if err := rows.Scan(&variable.Name, &variable.Value); err != nil {
			return nil, err
		}
		rouletteSettings = append(rouletteSettings, variable)
	}
	return rouletteSettings, nil
}
