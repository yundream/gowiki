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
	var ok bool
	session := sessions.SessionData{}
	cookie, _ := r.Cookie("session-jwt")
	fmt.Println("cookie", cookie.Value)
	if cookie != nil {
		session, ok = sessions.Validation(cookie.Value)
		if ok {
			return ("LOGIN User")
		}
	}

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("plugin/login/login.tmpl")
		if err != nil {
			return ""
		}

		a := struct{}{}
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

		fmt.Println(id, ":", password)
		q := c.Find(bson.M{"email": id, "password": enpass})
		n, err := q.Count()
		if err != nil {
			return err.Error()
		}
		if n == 0 {
			return "ID / Password FAIL"
		}
		if n == 1 {
			tokenString := sessions.Create(sessions.SessionData{id, "yundream@gmail.com", true})
			cookie := http.Cookie{Name: "session-jwt", Value: tokenString, Path: "/", Domain: "localhost"}
			http.SetCookie(w, &cookie)
			return "Success"
		}
	}
	return "FAIL"
}
