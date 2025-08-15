package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/architecture/web/handler/cookies"
	"forum/architecture/web/handler/view"

	ssession "forum/architecture/service/session"
	suser "forum/architecture/service/user"
)

// SignInHandler -
func (m *MainHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	debugLogHandler("SignInHandler", r)

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
			log.Printf("SignInHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
	}

	// Logic
	switch r.Method {
	case http.MethodGet:
		cookies.AddRedirectCookie(w, r.URL.Query().Get("redirect_to"))
		m.view.ExecuteTemplate(w, nil, "sign-in.html")
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("SignInHandler: r.ParseForm: %v\n", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		}

		usr, err := m.service.User.GetByNicknameOrEmail(r.FormValue("login"))
		switch {
		case err == nil:
		case errors.Is(err, suser.ErrNotFound):
			pg := &view.Page{Error: fmt.Errorf("user with login \"%v\" not found", r.FormValue("login"))}
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		case errors.Is(err, suser.ErrInvalidEmail):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid email %v", r.FormValue("login"))}
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		case errors.Is(err, suser.ErrInvalidNickname):
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid nickname %v", r.FormValue("login"))}
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		default:
			log.Printf("SignInHandler: User.GetByNicknameOrEmail: %s", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		}

		areEqual, err := usr.CompareHashAndPassword(r.FormValue("password"))
		switch {
		case err != nil:
			log.Printf("SignInHandler: user.CompareHashAndPassword: %s", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		case !areEqual:
			w.WriteHeader(http.StatusBadRequest)
			pg := &view.Page{Error: fmt.Errorf("invalid password for login \"%s\"", r.FormValue("login"))}
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		}

		session, err := m.service.Session.Record(usr.Id)
		if err != nil {
			log.Printf("SignInHandler: Session.Record: %s", err)
			pg := &view.Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.view.ExecuteTemplate(w, pg, "sign-in.html")
			return
		}
		expiresAfterSeconds := time.Until(session.ExpiredAt).Seconds()
		cookies.AddSessionCookie(w, session.Uuid, int(expiresAfterSeconds))

		if cookie := cookies.GetRedirectCookie(w, r); cookie != nil {
			cookies.RemoveRedirectCookie(w, r)
			http.Redirect(w, r, cookie.Value, http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
