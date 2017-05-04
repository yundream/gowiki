package main

import (
	"flag"
	"github.com/yundream/gowiki/app"
	"os"
)

const (
	RUN_MODE_CFG     = "/opt/joinc/gowiki/cfg/config.ini"
	INSTALL_MODE_CFG = "/opt/joinc/gowiki/cfg/install.ini"
)

func main() {
	cfgFile := flag.String("config", "/opt/joinc/gowiki/cfg/config.ini", "-config=config.ini")
	flag.Parse()

	mode := app.MOD_RUN

	if _, err := os.Stat("/opt/joinc/gowiki/cfg/config.ini"); os.IsNotExist(err) {
		mode = app.MOD_INSTALL
	}

	var wiki *app.Application
	switch mode {
	case app.MOD_RUN:
		wiki = app.New(*cfgFile)
		wiki.Run()
	case app.MOD_INSTALL:
		wiki = app.New(INSTALL_MODE_CFG)
		wiki.Run()
	}
}
