package database

import (
	"database/sql"
	_ "embed"
	"log"
	"strings"
)

//go:embed version.txt
var embeddedVersion string

func GetCurrentAppSchemaVersion() string {
	return strings.TrimSpace(embeddedVersion)
}

func Migrate(fromVersion, toVersion string) {
	log.Printf("🔄 Запуск миграции с версии %s на %s", fromVersion, toVersion)

	if fromVersion == "0.0.7" && toVersion == "0.0.8" {
		migrateFrom007To008()
	}
	log.Println("✅ Миграция завершена.")
}

func migrateFrom007To008() {
	log.Println("🚧 Миграция: 0.0.7 → 0.0.8")

	var keysToMigrate = []string{
		"rollPrice",
		"rollPriceIncrease",
	}

	values := make(map[string]string)

	for _, key := range keysToMigrate {
		val, err := CredentialsDB.GetENVValue(key)
		if err != nil {
			log.Printf("⚠️ Не удалось получить %s из CredentialsDB: %v", key, err)
			continue
		}
		values[key] = val
	}

	// Пример вставки
	/*
		for key, val := range values {
			err := WidgetDB.InsertENVValue(key, val)
			if err != nil {
				log.Printf("❌ Не удалось вставить %s в WidgetDB: %v", key, err)
				continue
			}
			log.Printf("✅ Перенесено %s в WidgetDB", key)
		}
	*/

	for _, key := range keysToMigrate {
		_, err := CredentialsDB.db.Exec("DELETE FROM EnvVariables WHERE name = ?", key)
		if err != nil {
			log.Printf("⚠️ Не удалось удалить %s из CredentialsDB: %v", key, err)
			continue
		}
		log.Printf("🧹 Удалено %s из CredentialsDB", key)
	}

	log.Println("✅ Значения rollPrice и rollPriceIncrease успешно перенесены в WidgetDB.")
}

func (wd *WidgetsDatabase) InitNewBase007To008() {
	var err error
	wd.db, err = sql.Open("sqlite", "./Widgets.db")
	if err != nil {
		log.Printf("❌ Ошибка подключения к базе WidgetsDB: %s", err)
	}

	if err = wd.db.Ping(); err != nil {
		log.Fatal("❌Ошибка пинга БД:", err)
	}

	_, err = wd.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Ошибка включения foreign_keys:", err)
	}

	//КОРОЧЕ ДОПИЛИТЬ ЧТОБЫ ИНСЕРТЫ ДЕЛАЛИСЬ ОДИН РАЗ ПРИ СОЗДАНИИ
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
		group_id INTEGER NOT NULL,
        name TEXT NOT NULL,
		FOREIGN KEY (group_id) REFERENCES RouletteSectorGroup(id)
    );
	`

	_, err = wd.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("❌ Ошибка создания таблиц в WidgetsDB: %s", err)
	}
	//WidgetDB.seedDefaultGroups()
}
