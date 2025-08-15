package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forum/architecture/models"
	"forum/architecture/service/post_comment"
	suser "forum/architecture/service/user"
	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"
)

// PostCommentCreateHandler -
func (m *MainHandler) PostCommentCreateHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostCommentCreateHandler", r)

	// Allowed Methods
	switch r.Method {
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCommentCreateHandler: r.Context().Value(\"UserId\") is nil")
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "post-create.html")
		return
	}

	userId := iUserId.(int64)
	user, err := m.service.User.GetByID(userId)
	switch {
	case err == nil:
	case errors.Is(err, suser.ErrNotFound):
		cookies.RemoveSessionCookie(w, r)
		if http.MethodGet == r.Method {
			cookies.AddRedirectCookie(w, r.RequestURI)
		}
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	case err != nil:
		log.Printf("PostEditHandler: m.service.User.GetByID: %v\n", err)
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
		return
	}

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()

		strPostId := r.FormValue("post_id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		comment := &models.PostComment{
			Content: r.FormValue("content"),
			PostId:  postId,
			UserId:  user.Id,
		}
		_, err = m.service.PostComment.Create(comment)
		switch {
		case err == nil:
		case errors.Is(err, post_comment.ErrInvalidContentLength):
			http.Error(w, "invalid length of content", http.StatusBadRequest)
			return
		default:
			log.Printf("PostCommentCreateHandler: m.service.PostComment.Create: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
