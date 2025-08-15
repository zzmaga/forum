package models

import "time"

// User представляет пользователя форума
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Пароль не отправляется в JSON
}

// Session представляет сессию пользователя
type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Category представляет категорию постов
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Post представляет пост на форуме
type Post struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	UserID       int64      `json:"user_id"`
	Username     string     `json:"username"`
	CreatedAt    string     `json:"created_at"`
	CommentCount int        `json:"comment_count"`
	Likes        int        `json:"likes"`
	Dislikes     int        `json:"dislikes"`
	Categories   []Category `json:"categories,omitempty"`
}

// Comment представляет комментарий к посту
type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	PostID    int64  `json:"post_id"`
	CreatedAt string `json:"created_at"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
}

// PostCategory связывает посты и категории (many-to-many)
type PostCategory struct {
	ID         int64 `json:"id"`
	PostID     int64 `json:"post_id"`
	CategoryID int64 `json:"category_id"`
}

// PostVote представляет голос за пост (лайк/дизлайк)
type PostVote struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"user_id"`
	PostID   int64 `json:"post_id"`
	IsLike   bool  `json:"is_like"`
}

// PostCommentVote представляет голос за комментарий (лайк/дизлайк)
type PostCommentVote struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	CommentID int64 `json:"comment_id"`
	IsLike    bool  `json:"is_like"`
} 