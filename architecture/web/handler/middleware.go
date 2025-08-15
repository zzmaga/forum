package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"forum/architecture/web/handler/cookies"

	ssession "forum/architecture/service/session"
)

func (m *MainHandler) MiddlewareMethodChecker(next http.Handler, allowedMthods map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		debugLogHandler("MiddlewareMethodChecker", r)
		if _, ok := allowedMthods[r.Method]; !ok {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareSessionChecker - NOT FINISHED
func (m *MainHandler) MiddlewareSessionChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		debugLogHandler("MiddlewareSessionChecker", r)
		cookie := cookies.GetSessionCookie(w, r)
		if cookie == nil {
			if r.Method == http.MethodGet {
				cookies.AddRedirectCookie(w, r.RequestURI)
			}
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		session, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
		case errors.Is(err, ssession.ErrExpired) || errors.Is(err, ssession.ErrNotFound):
			if r.Method == http.MethodGet {
				cookies.AddRedirectCookie(w, r.RequestURI)
			}
			cookies.RemoveSessionCookie(w, r)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		case err != nil:
			log.Printf("MiddlewareSessionChecker: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "UserId", session.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
