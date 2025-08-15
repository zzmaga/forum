package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"forum/architecture/models"
	"forum/architecture/service/post_comment_vote"
)

// PostCommentVoteHandler -
func (m *MainHandler) PostCommentVoteHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostCommentVoteHandler", r)

	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCommentVoteHandler: r.Context().Value(\"UserId\") is nil")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId := iUserId.(int64)

	switch r.Method {
	case http.MethodGet:
		strCommentId := r.URL.Query().Get("comment_id")
		commentId, err := strconv.ParseInt(strCommentId, 10, 64)
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

		postVote := &models.PostCommentVote{CommentId: commentId, UserId: userId, Vote: int8(vote)}
		err = m.service.PostCommentVote.Record(postVote)
		switch {
		case err == nil:
		case errors.Is(err, post_comment_vote.ErrInvalidVote) || errors.Is(err, post_comment_vote.ErrNotFound):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		case err != nil:
			log.Printf("PostCommentVoteHandler: m.service.PostCommentVote.Record: %s", err)
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
