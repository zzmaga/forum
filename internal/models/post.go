package models

type Post struct {
	ID      int
	Title   string
	Content string
	Author  string
}

var posts []Post

func AddPost(post Post) {
	posts = append(posts, post)
}

func GetPosts() []Post {
	return posts
}

func DeletePost(id int) {
	for i, post := range posts {
		if post.ID == id {
			posts = append(posts[:i], posts[i+1:]...)
			break
		}
	}
}
