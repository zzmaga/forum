package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"

	suser "forum/architecture/service/user"
)

// PostsOwnHandler -
func (m *MainHandler) PostsOwnHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostsOwnHandler", r)

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostsOwnHandler: r.Context().Value(\"UserId\") is nil")
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
		return
	}

	userId := iUserId.(int64)
	user, err := m.service.User.GetByID(userId)
	switch {
	case err == nil:
	case errors.Is(err, suser.ErrNotFound):
		cookies.RemoveSessionCookie(w, r)
		cookies.AddRedirectCookie(w, r.RequestURI)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	case err != nil:
		log.Printf("PostsOwnHandler: m.service.User.GetByID: %v\n", err)
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
		return
	}

	switch r.Method {
	case http.MethodGet:
		posts, err := m.service.Post.GetByUserID(user.Id, 0, 0)
		if err != nil {
			log.Printf("PostsOwnHandler: GetByUserID: %v\n", err)
		}

		err = m.service.FillPosts(posts, user.Id)
		if err != nil {
			log.Printf("PostsOwnHandler: FillPosts: %v\n", err)
		}

		pg := &view.Page{User: user, Posts: posts, Info: fmt.Errorf("Here is your posts")}
		m.view.ExecuteTemplate(w, pg, "home.html")
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
