package handlers

import (
	"database/sql"
	"forum/internal/database"
	"forum/internal/models"
	internal "forum/internal/template"
	"log"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCreatePostHandler(w, r)
	case http.MethodPost:
		PostCreatePostHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	categories, err := database.GetCategories()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Categories": categories,
		"UserID":     userID,
	}

	internal.RenderTemplate(w, "create_post.html", data)
}

func PostCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.Form.Get("title")
	content := r.Form.Get("content")
	categories := r.Form["categories"]

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	postID, err := database.CreatePost(userID, title, content)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Добавляем категории к посту
	for _, catIDStr := range categories {
		if catIDStr != "" {
			catID, catErr := strconv.Atoi(catIDStr)
			if catErr == nil {
				catErr = database.AddCategoryToPost(int(postID), catID)
				if catErr != nil {
					log.Printf("Warning: failed to add category %s to post: %v", catIDStr, catErr)
				}
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		http.Error(w, "Post ID required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Получаем пост по ID
	targetPost, err := database.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Получаем комментарии
	comments, err := database.GetCommentsByPost(postID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Post":     targetPost,
		"Comments": comments,
	}

	internal.RenderTemplate(w, "view_post.html", data)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := GetUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	postIDStr := r.Form.Get("post_id")
	content := r.Form.Get("content")

	if postIDStr == "" || content == "" {
		http.Error(w, "Post ID and content are required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Проверяем, что пост существует
	_, err = database.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	err = database.CreateComment(userID, postID, content)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+postIDStr, http.StatusSeeOther)
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	postIDStr := r.Form.Get("post_id")
	commentIDStr := r.Form.Get("comment_id")
	isLikeStr := r.Form.Get("is_like")

	var postID, commentID int
	var isLike bool

	if postIDStr != "" {
		var err error
		postID, err = strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Проверяем, что пост существует
		_, err = database.GetPostByID(postID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}
	}

	if commentIDStr != "" {
		var err error
		commentID, err = strconv.Atoi(commentIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Проверяем, что комментарий существует
		_, err = database.GetCommentByID(commentID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Comment not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}
	}

	if isLikeStr == "1" {
		isLike = true
	} else if isLikeStr == "0" {
		isLike = false
	} else {
		http.Error(w, "Invalid like value", http.StatusBadRequest)
		return
	}

	err = database.ToggleLike(userID, postID, commentID, isLike)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

func FilterPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filterType := r.URL.Query().Get("type")
	var posts []models.Post
	var err error

	switch filterType {
	case "category":
		categoryIDStr := r.URL.Query().Get("category_id")
		if categoryIDStr == "" {
			http.Error(w, "Category ID required", http.StatusBadRequest)
			return
		}
		var categoryID int
		categoryID, err = strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		posts, err = database.GetPostsByCategory(categoryID)

	case "user":
		userID, userErr := GetUserIDFromSession(r)
		if userErr != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		posts, err = database.GetPostsByUser(userID)

	case "liked":
		userID, userErr := GetUserIDFromSession(r)
		if userErr != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		posts, err = database.GetLikedPostsByUser(userID)

	default:
		posts, err = database.GetPosts()
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	categories, err := database.GetCategories()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Posts":      posts,
		"Categories": categories,
		"FilterType": filterType,
	}

	internal.RenderTemplate(w, "index.html", data)
}
