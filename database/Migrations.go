package database

import (
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
			err := RouletteDB.InsertENVValue(key, val)
			if err != nil {
				log.Printf("❌ Не удалось вставить %s в RouletteDB: %v", key, err)
				continue
			}
			log.Printf("✅ Перенесено %s в RouletteDB", key)
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

	log.Println("✅ Значения rollPrice и rollPriceIncrease успешно перенесены в RouletteDB.")
}
