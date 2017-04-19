package main

import (
	"github.com/yundream/gowiki/sessions"
	"net/http"
	"time"
)

func Function_logout(address string, parameter string, w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie("session-jwt")
	if cookie != nil {
		_, ok := sessions.Validation(cookie.Value)
		if ok {
			deleteCookie := http.Cookie{Name: "session-jwt",
				Value:   "none",
				Path:    "/",
				Domain:  "localhost",
				Expires: time.Now()}
			http.SetCookie(w, &deleteCookie)
			http.Redirect(w, r, "/w/login", http.StatusMovedPermanently)
			return ""
		}
	}
	return ""
}
