package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"forum/architecture/models"
	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"

	ssession "forum/architecture/service/session"
	suser "forum/architecture/service/user"
)

// SignUpHandler -
func (m *MainHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("SignUpHandler", r)

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cookie := cookies.GetSessionCookie(w, r)
	switch {
	case cookie == nil:
	case cookie != nil:
		_, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
			http.Redirect(w, r, "/", http.StatusFound)
			return
		case errors.Is(err, ssession.ErrExpired) || errors.Is(err, ssession.ErrNotFound):
			cookies.AddRedirectCookie(w, r.RequestURI)
			cookies.RemoveSessionCookie(w, r)
		case err != nil:
			log.Printf("SignUpHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
	}

	// Logic
	switch r.Method {
	case http.MethodGet:
		cookies.AddRedirectCookie(w, r.URL.Query().Get("redirect_to"))
		m.view.ExecuteTemplate(w, nil, "sign-up.html")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("SignUpHandler: r.ParseForm: %v\n", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		}

		newUser := &models.User{
			Nickname: r.FormValue("nickname"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		_, err = m.service.User.Create(newUser)
		switch {
		case err == nil:
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		case errors.Is(err, suser.ErrExistNickname):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("nickname \"%v\" is used. Try with another nickname.", newUser.Nickname)}
			m.view.ExecuteTemplate(w, pg, "sign-up.html")
		case errors.Is(err, suser.ErrExistEmail):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("email \"%v\" is used. Try with another email.", newUser.Email)}
			m.view.ExecuteTemplate(w, pg, "sign-up.html")
		case errors.Is(err, suser.ErrInvalidNickname):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid nickname \"%v\"", newUser.Nickname)}
			m.view.ExecuteTemplate(w, pg, "sign-up.html")
		case errors.Is(err, suser.ErrInvalidEmail):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid email \"%v\"", newUser.Email)}
			m.view.ExecuteTemplate(w, pg, "sign-up.html")
		default:
			log.Printf("SignUpHandler: %s", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "sign-up.html")
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
