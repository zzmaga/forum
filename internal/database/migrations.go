package database

import "log"

// -- USERS
func UsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		username NVARCHAR(32) UNIQUE NOT NULL CHECK(LENGTH(username) <= 32),
		email NVARCHAR(320) UNIQUE NOT NULL CHECK(LENGTH(email) <= 320),
		password TEXT
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}
}

// -- SESSIONS
func SessionsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY NOT NULL,
		user_id INTEGER NOT NULL,
		expired_at TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating sessions table:", err)
	}
}

// -- CATEGORIES
func CategoriesTable() {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating categories table:", err)
	}
}

// -- POSTS
func PostsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		title NVARCHAR(100) NOT NULL CHECK(LENGTH(title) <= 100),
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating posts table:", err)
	}
}

// -- POST CATEGORIES (many-to-many)
func PostCategoriesTable() {
	query := `
	CREATE TABLE IF NOT EXISTS post_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		category_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		UNIQUE (category_id, post_id),
		FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating post_categories table:", err)
	}
}

// -- LIKES (универсальная таблица для лайков постов и комментариев)
func PostsLikesTable() {
	query := `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER,
		comment_id INTEGER,
		is_like INTEGER NOT NULL CHECK(is_like IN(0, 1)),
		UNIQUE (user_id, post_id, comment_id),
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
		CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating likes table:", err)
	}
}

// -- COMMENTS
func CommentsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating comments table:", err)
	}
}

// -- LIKES ON COMMENTS (эта функция больше не нужна, так как используем универсальную таблицу likes)
func CommentsLikesTable() {
	// Эта функция оставлена для совместимости, но не создает отдельную таблицу
	// Все лайки теперь хранятся в таблице likes
}
