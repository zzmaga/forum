package database

import "log"

// -- USERS
func UsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		nickname NVARCHAR(32) UNIQUE NOT NULL CHECK(LENGTH(nickname) <= 32),
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
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		uuid TEXT NOT NULL,
		expired_at TEXT,
		user_id INTEGER NOT NULL UNIQUE,
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
		name TEXT UNIQUE NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		description TEXT
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
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
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
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		UNIQUE (category_id, post_id),
		FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating post_categories table:", err)
	}
}

// -- POSTS LIKES
func PostsLikesTable() {
	query := `
	CREATE TABLE IF NOT EXISTS posts_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		vote INTEGER NOT NULL CHECK(vote IN(-1, 0, 1)),
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		UNIQUE (user_id, post_id),
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating posts_likes table:", err)
	}
}

// -- COMMENTS
func CommentsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS posts_comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating comments table:", err)
	}
}

// -- LIKES ON COMMENTS
func CommentsLikesTable() {
	query := `
	CREATE TABLE IF NOT EXISTS comments_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		vote INTEGER NOT NULL CHECK(vote IN(-1, 0, 1)),
		user_id INTEGER NOT NULL,
		comment_id INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		UNIQUE (user_id, comment_id),
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(comment_id) REFERENCES posts_comments(id) ON DELETE CASCADE
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating comments_likes table:", err)
	}
}
