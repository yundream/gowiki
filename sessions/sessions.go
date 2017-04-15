package sessions

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	MySecret = "TPW9EhIVuXCTAnLfcdZvup3JkLatqvnv"
)

type SessionData struct {
	Name  string
	Email string
	Admin bool
	Login bool
}

type UserClaims struct {
	*jwt.StandardClaims
	Level   string
	Session SessionData
}

func Create(sess SessionData) string {
	ttl := time.Now().Add(time.Hour * 24).Unix()
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &UserClaims{
		&jwt.StandardClaims{
			ExpiresAt: ttl,
		},
		"level1",
		sess,
	}
	tokenString, err := t.SignedString([]byte(MySecret))
	if err != nil {
		fmt.Println(err.Error())
	}
	return tokenString
}

func Validation(tokenString string) (SessionData, bool) {
	Claims := &UserClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		Claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(MySecret), nil
		})
	if err != nil {
		return SessionData{Login: false}, false
	}
	Claims.Session.Login = true
	return Claims.Session, true
}
