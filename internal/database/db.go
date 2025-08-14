package database

import (
	"database/sql"
	"forum/internal/models"
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

	UsersTable()
	SessionsTable()
	CategoriesTable()
	PostsTable()
	PostCategoriesTable()
	PostsLikesTable()
	CommentsTable()
	CommentsLikesTable()
	insertDefaultCategories()
}

// Добавляем базовые категории
func insertDefaultCategories() {
	categories := []string{"Общие", "Технологии", "Наука", "Искусство", "Спорт", "Политика"}
	for _, cat := range categories {
		_, err := DB.Exec("INSERT OR IGNORE INTO categories(name) VALUES (?)", cat)
		if err != nil {
			log.Printf("Warning: failed to insert category %s: %v", cat, err)
		}
	}
}

// Функции для работы с постами
func GetPostByID(postID int) (*models.Post, error) {
	row := DB.QueryRow(`
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
		       COUNT(DISTINCT c.id) as comment_count,
		       COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id AND l.comment_id IS NULL
		WHERE p.id = ?
		GROUP BY p.id
	`, postID)

	var p models.Post
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Username, &p.CommentCount, &p.Likes, &p.Dislikes)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func CreatePost(userID int, title, content string) (int64, error) {
	result, err := DB.Exec("INSERT INTO posts(user_id, title, content, created_at) VALUES (?, ?, ?, datetime('now'))", userID, title, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetPosts() ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
		       COUNT(DISTINCT c.id) as comment_count,
		       COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id AND l.comment_id IS NULL
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Username, &p.CommentCount, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostsByCategory(categoryID int) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username,
		       COUNT(DISTINCT c.id) as comment_count,
		       COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id AND l.comment_id IS NULL
		WHERE pc.category_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Username, &p.CommentCount, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostsByUser(userID int) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username,
		       COUNT(DISTINCT c.id) as comment_count,
		       COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id AND l.comment_id IS NULL
		WHERE p.user_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Username, &p.CommentCount, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetLikedPostsByUser(userID int) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username,
		       COUNT(DISTINCT c.id) as comment_count,
		       COALESCE(SUM(CASE WHEN l2.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l2.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN likes l ON p.id = l.post_id AND l.comment_id IS NULL
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l2 ON p.id = l2.post_id AND l2.comment_id IS NULL
		WHERE l.user_id = ? AND l.is_like = 1
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Username, &p.CommentCount, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// Функции для работы с комментариями
func CreateComment(userID, postID int, content string) error {
	_, err := DB.Exec("INSERT INTO comments(user_id, post_id, content, created_at) VALUES (?, ?, ?, datetime('now'))", userID, postID, content)
	return err
}

func GetCommentByID(commentID int) (*models.Comment, error) {
	row := DB.QueryRow(`
		SELECT c.id, c.content, c.created_at, u.username,
		       COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN likes l ON c.id = l.comment_id AND l.post_id IS NULL
		WHERE c.id = ?
		GROUP BY c.id
	`, commentID)

	var c models.Comment
	err := row.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.Username, &c.Likes, &c.Dislikes)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetCommentsByPost(postID int) ([]models.Comment, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.content, c.created_at, u.username,
		       COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN likes l ON c.id = l.comment_id AND l.post_id IS NULL
		WHERE c.post_id = ?
		GROUP BY c.id
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.Username, &c.Likes, &c.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// Функции для работы с лайками
func ToggleLike(userID, postID, commentID int, isLike bool) error {
	// Сначала удаляем существующий лайк/дизлайк
	_, err := DB.Exec("DELETE FROM likes WHERE user_id = ? AND ((post_id = ? AND comment_id IS NULL) OR (comment_id = ? AND post_id IS NULL))", userID, postID, commentID)
	if err != nil {
		return err
	}

	// Добавляем новый лайк/дизлайк
	var isLikeInt int
	if isLike {
		isLikeInt = 1
	} else {
		isLikeInt = 0
	}

	if postID > 0 {
		_, err = DB.Exec("INSERT INTO likes(user_id, post_id, comment_id, is_like) VALUES (?, ?, NULL, ?)", userID, postID, isLikeInt)
	} else {
		_, err = DB.Exec("INSERT INTO likes(user_id, post_id, comment_id, is_like) VALUES (?, NULL, ?, ?)", userID, commentID, isLikeInt)
	}
	return err
}

// Функции для работы с категориями
func GetCategories() ([]models.Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func AddCategoryToPost(postID, categoryID int) error {
	_, err := DB.Exec("INSERT INTO post_categories(post_id, category_id) VALUES (?, ?)", postID, categoryID)
	return err
}
