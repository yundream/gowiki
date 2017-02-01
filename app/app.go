package app

import (
	"github.com/yundream/gowiki/handler"
	"log"
)

type Application struct {
	Port string
}

func New(port string) *Application {
	return &Application{Port: port}
}

func (a Application) Run() {
	h := handler.New()
	h.DocumentDirectory = "/opt/gowiki/doc"
	h.Theme = "joinc"
	err := h.Run(a.Port)
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
}
