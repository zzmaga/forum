package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"forum/architecture/models"
	spost_vote "forum/architecture/service/post_vote"
)

// PostVoteHandler -
func (m *MainHandler) PostVoteHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostVoteHandler", r)

	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostVoteHandler: r.Context().Value(\"UserId\") is nil")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId := iUserId.(int64)

	switch r.Method {
	case http.MethodGet:
		strPostId := r.URL.Query().Get("post_id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		strVote := r.URL.Query().Get("vote")
		vote, err := strconv.ParseInt(strVote, 10, 8)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		postVote := &models.PostVote{PostId: postId, UserId: userId, Vote: int8(vote)}
		err = m.service.PostVote.Record(postVote)
		switch {
		case err == nil:
		case errors.Is(err, spost_vote.ErrInvalidVote) || errors.Is(err, spost_vote.ErrNotFound):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		case err != nil:
			log.Printf("PostVoteHandler: m.service.PostVote.Record: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
