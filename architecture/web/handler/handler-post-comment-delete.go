package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	spost "forum/architecture/service/post"
	"forum/architecture/web/handler/view"
)

// PostCommentDeleteHandler -
func (m *MainHandler) PostCommentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostCommentDeleteHandler", r)

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCommentDeleteHandler: r.Context().Value(\"UserId\") is nil")
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
		return
	}

	userId := iUserId.(int64)

	switch r.Method {
	case http.MethodGet:
		strPostCommentId := r.URL.Query().Get("id")
		postCommentId, err := strconv.ParseInt(strPostCommentId, 10, 64)
		if err != nil || postCommentId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		comment, err := m.service.PostComment.GetByID(postCommentId)
		switch {
		case err == nil:
		case errors.Is(err, spost.ErrNotFound):
			// TODO: error page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if comment.UserId != userId {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		err = m.service.PostComment.DeleteByID(comment.Id)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("PostCommentDeleteHandler: m.service.PostComment.DeleteByID: %v\n", err)
			pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
