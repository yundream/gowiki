package app

import (
	"github.com/yundream/gowiki/handler"
	"gopkg.in/ini.v1"
	"log"
)

const (
	MOD_INSTALL = iota
	MOD_RUN
)

type Application struct {
	CfgFile string
	Mode    int
}

func New(cfgFile string) *Application {
	return &Application{CfgFile: cfgFile}
}

func (a Application) SetMode(mode int) {
	a.Mode = mode
}

func (a Application) Run() {
	cfg, err := ini.Load([]byte(""), a.CfgFile)
	if err != nil {
		panic(err)
	}
	servers, err := cfg.GetSection("server")
	if err != nil {
		panic(err)
	}
	address := servers.KeysHash()["address"]
	mongo := servers.KeysHash()["mongo"]

	h, err := handler.New(mongo)
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}

	err = h.LoadTemplate("joinc")
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
	err = h.Run(address)
	if err != nil {
		log.Fatal("Server run error : ", err.Error())
	}
}

func (a Application) Installer() {

}
