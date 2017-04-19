package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/yundream/gowiki/sessions"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"net/http"
)

func Function_login(address string, parameter string, w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie("session-jwt")
	if cookie != nil {
		info, ok := sessions.Validation(cookie.Value)
		if ok {
			fmt.Println(info)
			t, err := template.ParseFiles("plugin/login/loginuser.tmpl")
			if err != nil {
				return ""
			}
			var doc bytes.Buffer
			err = t.Execute(&doc, info)
			if err != nil {
				return err.Error()
			}
			return doc.String()
		}
	}

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("plugin/login/login.tmpl")
		if err != nil {
			return ""
		}
		reason := r.URL.Query().Get("t")
		if len(reason) != 0 {
		}

		a := struct{ T string }{T: reason}
		var doc bytes.Buffer
		err = t.Execute(&doc, a)

		if err != nil {
			return err.Error()
		}
		return doc.String()
	case "POST":
		s, err := mgo.Dial(address)
		if err != nil {
			return "error"
		}
		c := s.DB("user").C("user")
		id := r.PostFormValue("id")
		password := r.PostFormValue("password")
		h := sha256.New()
		io.WriteString(h, password)
		enpass := base64.StdEncoding.EncodeToString(h.Sum(nil))

		q := c.Find(bson.M{"email": id, "password": enpass})
		n, err := q.Count()
		if err != nil {
			return err.Error()
		}
		if n == 0 {
			http.Redirect(w, r, "/w/login?t=fail", http.StatusMovedPermanently)
		}
		if n == 1 {
			tokenString := sessions.Create(sessions.SessionData{id, "yundream@gmail.com", true, true})
			cookie := http.Cookie{Name: "session-jwt", Value: tokenString, Path: "/", Domain: "localhost"}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/w/login", http.StatusMovedPermanently)
			return ""
		}
	}
	return "FAIL"
}
