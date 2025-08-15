package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"

	suser "forum/architecture/service/user"
)

// PostsVotedHandler -
func (m *MainHandler) PostsVotedHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostsVotedHandler", r)

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostsVotedHandler: r.Context().Value(\"UserId\") is nil")
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
		log.Printf("PostsVotedHandler: m.service.User.GetByID: %v\n", err)
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
		return
	}

	switch r.Method {
	case http.MethodGet:
		strVote := r.URL.Query().Get("vote")
		vote, err := strconv.ParseInt(strVote, 10, 8)
		if err != nil || vote < -1 || 1 < vote {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		postIDs, err := m.service.PostVote.GetAllUserVotedPostIDs(userId, int8(vote), 0, 0)
		if err != nil {
			log.Printf("PostsVotedHandler: PostVote.GetAllUserVotedPostIDs: %v\n", err)
			pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
			return
		}

		posts, err := m.service.Post.GetByIDs(postIDs)
		if err != nil {
			log.Printf("PostsVotedHandler: Post.GetByIDs: %v\n", err)
		}

		err = m.service.FillPosts(posts, user.Id)
		if err != nil {
			log.Printf("PostsVotedHandler: FillPosts: %v\n", err)
		}

		pg := &view.Page{User: user, Posts: posts, Info: fmt.Errorf("Voted Posts")}
		m.view.ExecuteTemplate(w, pg, "home.html")
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
