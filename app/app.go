package app

import (
	"github.com/yundream/gowiki/handler"
	"log"
	"net/http"
)

type Application struct {
	Port string
}

func New(port string) *Application {
	return &Application{Port: port}
}

func (a Application) Run() {
	h := handler.New()
	http.Handle("/", h.Router)
	err := http.ListenAndServe(a.Port, nil)
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
}
