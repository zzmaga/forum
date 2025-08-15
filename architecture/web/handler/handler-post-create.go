package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"forum/architecture/models"
	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"

	scategory "forum/architecture/service/category"
	spost "forum/architecture/service/post"
	suser "forum/architecture/service/user"
)

// PostCreateHandler -
func (m *MainHandler) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("PostCreateHandler", r)

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
		log.Println("PostCreateHandler: r.Context().Value(\"UserId\") is nil")
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
		pg := &view.Page{User: user}
		m.view.ExecuteTemplate(w, pg, "post-create.html")
		return
	case http.MethodPost:
		r.ParseForm()

		post := &models.Post{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			UserId:  userId,
		}
		_, err := m.service.Post.Create(post)
		switch {
		case err == nil:
		case errors.Is(err, spost.ErrInvalidTitleLength):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid length of title")}
			m.view.ExecuteTemplate(w, pg, "post-create.html")
			return
		case errors.Is(err, spost.ErrInvalidContentLength):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid length of content")}
			m.view.ExecuteTemplate(w, pg, "post-create.html")
			return
		default:
			log.Printf("PostCreateHandler: m.service.Post.Create: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later: %s", err)}
			m.view.ExecuteTemplate(w, pg, "post-create.html")
			return
		}

		catNames := strings.Fields(r.Form.Get("categories"))
		err = m.service.Category.AddToPostByNames(catNames, post.Id)
		switch {
		case err == nil:
		case errors.Is(err, scategory.ErrCategoryLimitForPost):
			err = m.service.Post.DeleteByID(post.Id)
			if err != nil {
				log.Println("PostCreateHandler: m.service.Post.DeleteByID: %w", err)
			}

			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("post not created, invalid categies count, category limit = %v", models.MaxCategoryLimitForPost)}
			m.view.ExecuteTemplate(w, pg, "post-create.html")
			return
		default:
			err = m.service.Post.DeleteByID(post.Id)
			if err != nil {
				log.Println("PostCreateHandler: m.service.Post.DeleteByID: %w", err)
			}

			log.Printf("PostCreateHandler:  m.service.Category.AddToPostByNames: %s", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later: %s", err)}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "post-create.html")
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/get?id=%v", post.Id), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
