package main

import (
	"flag"
	"github.com/yundream/gowiki/app"
)

func main() {
	port := flag.String("port", "0.0.0.0:8080", "-port=0.0.0.0:8080")
	flag.Parse()
	wiki := app.New(*port)
	wiki.Run()
}
