package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	createUsersTable()
}

func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}
}
