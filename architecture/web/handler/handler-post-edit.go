package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/architecture/models"
	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"

	scategory "forum/architecture/service/category"
	spost "forum/architecture/service/post"
	suser "forum/architecture/service/user"
)

// PostCreateHandler -
func (m *MainHandler) PostEditHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostEditHandler", r)

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostEditHandler: r.Context().Value(\"UserId\") is nil")
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "post-edit.html")
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
		log.Printf("PostEditHandler: m.service.User.GetByID: %v\n", err)
		pg := &view.Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.view.ExecuteTemplate(w, pg, "alert.html") // TODO: Custom Error Page
		return
	}

	switch r.Method {
	case http.MethodGet:
		strPostId := r.URL.Query().Get("id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}
		post, err := m.service.Post.GetByID(postId)
		switch {
		case err == nil:
		case errors.Is(err, spost.ErrNotFound):
			// TODO: error page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if post.UserId != user.Id {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		categories, err := m.service.Category.GetByPostID(post.Id)
		switch {
		case err == nil:
			post.WCategories = categories
		default:
			log.Printf("PostEditHandler: m.service.PostCategory.GetByPostID: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		pg := &view.Page{User: user, Post: post}
		m.view.ExecuteTemplate(w, pg, "post-edit.html")
		return
	case http.MethodPost:
		r.ParseForm()

		var strPostId string = r.FormValue("id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		post, err := m.service.Post.GetByID(postId)
		switch {
		case err == nil:
		case errors.Is(err, spost.ErrNotFound):
			// TODO: error page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if post.UserId != user.Id {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		post = &models.Post{
			Id:      postId,
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			UserId:  user.Id,
		}
		err = m.service.Post.Update(post)
		switch {
		case err == nil:
		case errors.Is(err, spost.ErrInvalidTitleLength) || errors.Is(err, spost.ErrInvalidContentLength):
			categories, errn := m.service.Category.GetByPostID(post.Id)
			switch {
			case errn == nil:
				post.WCategories = categories
			default:
				log.Printf("PostEditHandler: m.service.PostCategory.GetByPostID: %v\n", err)
				http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
				return
			}

			var errMsg error
			switch {
			case errors.Is(err, spost.ErrInvalidTitleLength):
				errMsg = fmt.Errorf("invalid title length")
			case errors.Is(err, spost.ErrInvalidContentLength):
				errMsg = fmt.Errorf("invalid content length")
			default:
				log.Printf("PostEditHandler: havent got message for error: %s\n", err)
				errMsg = fmt.Errorf("invalid post")
			}

			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{User: user, Post: post, Error: errMsg}
			m.view.ExecuteTemplate(w, pg, "post-edit.html")
			return
		case err != nil:
			log.Printf("PostEditHandler: m.service.Post.Update: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		err = m.service.Category.DeleteByPostID(post.Id)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("PostEditHandler: m.service.PostCategory.DeleteByPostID: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		catNames := strings.Fields(r.Form.Get("categories"))
		err = m.service.Category.AddToPostByNames(catNames, post.Id)
		switch {
		case err == nil:
		case errors.Is(err, scategory.ErrCategoryLimitForPost):
			pg := &view.Page{Warn: fmt.Errorf("post categories not updated, invalid categies count, category limit = %v", models.MaxCategoryLimitForPost), Post: post}
			w.WriteHeader(http.StatusBadRequest)
			m.view.ExecuteTemplate(w, pg, "post-edit.html")
			return
		default:
			log.Printf("PostEditHandler:  m.service.Category.AddToPostByNames: %s", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later: %s", err), Post: post}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "post-edit.html")
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/get?id=%v", post.Id), http.StatusSeeOther)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
