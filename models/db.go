package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB // глобальная переменная, доступна другим

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath) // ← не создаём новую переменную с :=
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		username TEXT,
		password TEXT
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
}

func RegisterUser(email, username, password string) error {
	_, err := DB.Exec(`
		INSERT INTO users (email, username, password)
		VALUES (?, ?, ?);
	`, email, username, password)
	return err
}
