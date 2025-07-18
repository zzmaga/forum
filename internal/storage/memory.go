package storage

import (
	"sync"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Author  string
}

type MemoryStorage struct {
	mu     sync.Mutex
	posts  []Post
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:  []Post{},
		nextID: 1,
	}
}

func (s *MemoryStorage) AddPost(title, content, author string) Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := Post{
		ID:      s.nextID,
		Title:   title,
		Content: content,
		Author:  author,
	}
	s.posts = append(s.posts, post)
	s.nextID++
	return post
}

func (s *MemoryStorage) GetPosts() []Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.posts
}

func (s *MemoryStorage) DeletePost(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, post := range s.posts {
		if post.ID == id {
			s.posts = append(s.posts[:i], s.posts[i+1:]...)
			return true
		}
	}
	return false
}
