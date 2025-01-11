package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB создаёт таблицу scheduler, если её нет
func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	log.Printf("✅🔌 Подключаемся к базе: %s", dbPath)

	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("❌ Ошибка при закрытии соединения с БД: %v", cerr)
		}
	}()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		date TEXT NOT NULL,
		title TEXT NOT NULL, 
		comment TEXT, 
		repeat TEXT(128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date); 
	CREATE INDEX IF NOT EXISTS idx_title ON scheduler(title);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("❌ Ошибка при создании таблицы: %v", err)
		return err
	}

	log.Printf("✅ Таблица scheduler в [%s] создана или уже есть", dbPath)
	return nil
}
