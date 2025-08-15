package cookies

import (
	"net/http"
)

const (
	CookieRedirectName = "redirect_to"
)

// AddRedirectCookie - sets redirect cookie if field redirectTo is not empty.
// field redirectTo sets at cookie Value
func AddRedirectCookie(w http.ResponseWriter, redirectTo string) {
	if redirectTo == "" {
		return
	}
	http.SetCookie(w,
		&http.Cookie{
			Name:   CookieRedirectName,
			Value:  redirectTo,
			Path:   "/",
			MaxAge: 3600,
		},
	)
}

// GetRedirectCookie - returns redirect cookie
func GetRedirectCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	return getCookieByName(w, r, CookieRedirectName)
}

// RemoveRedirectCookie - removes cookie by setting maxAge -1
func RemoveRedirectCookie(w http.ResponseWriter, r *http.Request) {
	removeCookieByName(w, r, CookieRedirectName)
}
