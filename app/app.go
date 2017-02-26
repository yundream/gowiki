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
	h, err := handler.New()
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
	h.DocumentDirectory = "/opt/gowiki/doc"
	err = h.LoadTemplate("joinc")
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
	err = h.Run(a.Port)
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
}
