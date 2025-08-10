package models

// Posts
type Post struct {
	ID           int
	Title        string
	Content      string
	CreatedAt    string
	Username     string
	CommentCount int
	Likes        int
	Dislikes     int
}

// Categories
type Category struct {
	ID   int
	Name string
}
