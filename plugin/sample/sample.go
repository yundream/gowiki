package main

import (
	"github.com/yundream/gowiki/handler"
)

func Function_sample(data string, opt handler.Options) string {
	return "Hello World " + opt.Name
}
