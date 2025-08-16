package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB opens sqlite database and runs architecture-compatible migrations
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}
	if err := insertDefaultCategoriesNew(db); err != nil {
		log.Printf("warning: insertDefaultCategories: %v", err)
	}

	// keep global for legacy parts if any
	DB = db
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		// users
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			nickname NVARCHAR(32) UNIQUE NOT NULL CHECK(LENGTH(nickname) <= 32),
			email NVARCHAR(320) UNIQUE NOT NULL CHECK(LENGTH(email) <= 320),
			password TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		// sessions
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			uuid TEXT UNIQUE NOT NULL,
			expired_at TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		// categories
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name TEXT UNIQUE NOT NULL CHECK(LENGTH(name) <= 64),
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		// posts
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			title NVARCHAR(100) NOT NULL CHECK(LENGTH(title) <= 100),
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		// posts_categories
		`CREATE TABLE IF NOT EXISTS posts_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			UNIQUE (post_id, category_id),
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
		);`,
		// posts_votes
		`CREATE TABLE IF NOT EXISTS posts_votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			vote INTEGER NOT NULL CHECK(vote IN(-1, 1)),
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user_id, post_id),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`,
		// posts_comments
		`CREATE TABLE IF NOT EXISTS posts_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`,
		// posts_comments_votes
		`CREATE TABLE IF NOT EXISTS posts_comments_votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			vote INTEGER NOT NULL CHECK(vote IN(-1, 1)),
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user_id, comment_id),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(comment_id) REFERENCES posts_comments(id) ON DELETE CASCADE
		);`,
	}
	for _, q := range stmts {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	// compatibility adjustments for legacy schemas
	// add updated_at to posts if missing
	ensureColumn(db, "posts", "updated_at", "TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP")
	// add created_at to categories if missing
	ensureColumn(db, "categories", "created_at", "TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP")
	// add uuid to sessions if missing
	ensureColumn(db, "sessions", "uuid", "TEXT")
	// add nickname to users if missing; if legacy username exists, copy values
	ensureColumn(db, "users", "nickname", "NVARCHAR(32)")
	_, _ = db.Exec(`UPDATE users SET nickname = email WHERE (nickname IS NULL OR nickname = '') AND email IS NOT NULL AND email != ''`)

	return nil
}

func ensureColumn(db *sql.DB, table, column, decl string) {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		cid     int
		name    string
		ctype   string
		notnull int
		dflt    interface{}
		pk      int
	)
	exists := false
	for rows.Next() {
		_ = rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		if name == column {
			exists = true
			break
		}
	}
	if !exists {
		_, _ = db.Exec("ALTER TABLE " + table + " ADD COLUMN " + column + " " + decl)
	}
}

func insertDefaultCategoriesNew(db *sql.DB) error {
	cats := []string{"General", "Technology", "Science", "Art", "Sports", "Politics"}
	for _, c := range cats {
		if _, err := db.Exec("INSERT OR IGNORE INTO categories(name) VALUES (?)", c); err != nil {
			return err
		}
	}
	return nil
}
